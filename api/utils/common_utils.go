package utils

import "api/models"

func ChunkSlicePosts(slice []models.Post, chunkSize int) [][]models.Post {
	var chunks [][]models.Post

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])

	}

	return chunks
}

func ChunkSlice(slice []any, chunkSize int) [][]any {
	var chunks [][]any

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])

	}

	return chunks
}
