package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"mime/multipart"
	"slices"
)

func uploadFiles(req httpc.RequestBuilder, data any, files []MessageFile) {
	isMultipart := slices.ContainsFunc(files, func(file MessageFile) bool {
		return file.Resolver != nil
	})
	if isMultipart {
		for i, file := range files {
			if file.Resolver != nil {
				file.ID = snowflake.ID(i)
			}
		}
		req.Multipart(func(writer *multipart.Writer) error {
			{
				buff, err := json.Marshal(data)
				if err != nil {
					return fmt.Errorf("failed to marshal payload: %w", err)
				}
				if err = writer.WriteField("payload_json", string(buff)); err != nil {
					return fmt.Errorf("failed to write payload to multipart.Writer: %w", err)
				}
			}
			for _, file := range files {
				if file.Resolver != nil {
					resolved, err := file.Resolver()
					if err != nil {
						return fmt.Errorf("failed to resolve attachment (%d - %s): %w", file.ID, file.Filename, err)
					}
					fileWriter, err := writer.CreateFormFile(fmt.Sprintf("files[%d]", file.ID), file.Filename)
					if err != nil {
						return fmt.Errorf("failed to create form file: %w", err)
					}
					_, err = fileWriter.Write(resolved)
					if err != nil {
						return fmt.Errorf("failed to write to form file: %w", err)
					}
				}
			}
			return nil
		})
	} else {
		req.Body(data)
	}
}
