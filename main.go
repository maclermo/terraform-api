package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"terraform-api/runner"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

func unzipSource(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf("error while opening zipped file from formData %s", err)
	}
	defer reader.Close()

	destination, err = filepath.Abs(destination)
	if err != nil {
		return fmt.Errorf("error while creating destination directory or directories %s", err)
	}

	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return fmt.Errorf("error while unzipping file or files %s", err)
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path %s", filePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return fmt.Errorf("destination is not a valid directory for writing %s", err)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return fmt.Errorf("cannot create directories for unzipped files %s", err)
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("cannot open destination file for writing %s", err)
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		return fmt.Errorf("cannot open zipped file %s", err)
	}

	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return fmt.Errorf("could not copy unzipped file to destination %s", err)
	}

	return nil
}

func saveFileHandler(c *gin.Context) (string, error) {
	file, err := c.FormFile("terraform")
	if err != nil {
		return "", fmt.Errorf("error while getting file %s", err)
	}

	extension := filepath.Ext(file.Filename)
	fileUuid := uuid.New()
	fileName := fileUuid.String() + extension
	zipPath := "/tmp/zip"
	zipFile := zipPath + "/" + fileName
	tfPath := "/tmp/tf/" + fileUuid.String()

	_ = os.Mkdir(zipPath, os.ModePerm)

	if err := c.SaveUploadedFile(file, zipFile); err != nil {
		return "", fmt.Errorf("error while saving file %s", err)
	}

	_ = os.Mkdir(tfPath, os.ModePerm)

	if err := unzipSource(zipFile, tfPath); err != nil {
		return "", fmt.Errorf("error while unzipping %s", err)
	}

	_ = os.Remove(zipFile)

	return tfPath, nil
}

func parseFormData(c *gin.Context) (runner.Request, error) {
	jsonRequest := runner.JSONRequest{}

	if err := c.MustBindWith(&jsonRequest, binding.Form); err != nil {
		return runner.Request{}, fmt.Errorf("binding failed with form %s", err)
	}

	vars := map[string]interface{}{}
	if err := json.Unmarshal([]byte(jsonRequest.Vars), &vars); err != nil && len(vars) > 0 {
		return runner.Request{}, fmt.Errorf("json unmarshal error with field vars %s", err)
	}

	envVars := map[string]string{}
	if err := json.Unmarshal([]byte(jsonRequest.EnvVars), &envVars); err != nil && len(envVars) > 0 {
		return runner.Request{}, fmt.Errorf("json unmarshal error with field envVars %s", err)
	}

	backendConfig := map[string]interface{}{}
	if err := json.Unmarshal([]byte(jsonRequest.BackendConfig), &backendConfig); err != nil && len(backendConfig) > 0 {
		return runner.Request{}, fmt.Errorf("json unmarshal error with field backendConfig %s", err)
	}

	request := runner.Request{
		Workspace:     jsonRequest.Workspace,
		Vars:          vars,
		EnvVars:       envVars,
		BackendConfig: backendConfig,
	}

	return request, nil
}

func dispatchActions(c *gin.Context) (string, runner.Request, error) {
	tfPath, err := saveFileHandler(c)
	if err != nil {
		return "", runner.Request{}, fmt.Errorf("error parsing saveFileHandler %s", err)
	}

	requestBind, err := parseFormData(c)
	if err != nil {
		return "", runner.Request{}, fmt.Errorf("error parsing parseFormData %s", err)
	}

	return tfPath, requestBind, nil
}

func main() {
	runner.InitJobs()

	router := gin.Default()

	router.POST("/plan", func(c *gin.Context) {
		tfPath, requestBind, err := dispatchActions(c)
		if err != nil {
			c.JSON(500, map[string]string{"error": fmt.Sprint(err)})
			return
		}
		response := plan(tfPath, requestBind)
		c.JSON(200, response)
	})

	router.POST("/apply", func(c *gin.Context) {
		tfPath, requestBind, err := dispatchActions(c)
		if err != nil {
			c.JSON(500, map[string]string{"error": fmt.Sprint(err)})
			return
		}
		response := apply(tfPath, requestBind)
		c.JSON(200, response)
	})

	router.POST("/destroy", func(c *gin.Context) {
		tfPath, requestBind, err := dispatchActions(c)
		if err != nil {
			c.JSON(500, map[string]string{"error": fmt.Sprint(err)})
			return
		}
		response := destroy(tfPath, requestBind)
		c.JSON(200, response)
	})

	router.GET("/output/:id", func(c *gin.Context) {
		output, err := output(c.Param("id"))
		if err != nil {
			c.JSON(500, map[string]string{"error": fmt.Sprint(err)})
			return
		}
		c.JSON(200, output)
	})

	router.GET("/status/:id", func(c *gin.Context) {
		status, err := status(c.Param("id"))
		if err != nil {
			c.JSON(500, map[string]string{"error": fmt.Sprint(err)})
			return
		}
		c.JSON(200, status)
	})

	router.GET("/delete/:id", func(c *gin.Context) {
		deleted := delete(c.Param("id"))
		c.JSON(200, deleted)
	})

	router.GET("/result/:id", func(c *gin.Context) {
		result, err := result(c.Param("id"))
		if err != nil {
			c.JSON(500, map[string]string{"error": fmt.Sprint(err)})
		}
		c.JSON(200, result)
	})

	router.Run()
}
