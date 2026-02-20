// Copyright IBM Corp. 2024
// SPDX-License-Identifier: MPL-2.0

package generator

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	openai "github.com/sashabaranov/go-openai"
)

const fileHeader = `package %s

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)
`

const systemPrompt = `You are a Terraform Provider Developer. You are writing code for the Kubernetes Provider using the newer terraform-plugin-framework. You will receive an attribute description and you will write a single Validator implementation that meets all the restrictions detailed in the description. The beginning of the description will provide the base type def for you to use. Do not perform any explanations, just write code. Do not define the whole file with imports or a main function, just write the interface implementation. If a package is to be imported do not assume an alias name use the package name directly, example: github.com/hashicorp/terraform-plugin-framework/schema/validator use "validator" but remember don't write the code imports. Be careful of the implied type, for a string validator the function to implement is ValidateString etc. Make sure you define the type struct that is passed to you. Make sure Description MarkdownDescription are implemented. There is no Validate function to implement, only Validate[Type].`

var client *openai.Client

func generateValidator(contents string, attrs AttributesGenerator) string {
	if client == nil {
		client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	}

	for _, attr := range attrs {

		if len(attr.NestedAttributes) != 0 {
			contents = generateValidator(contents, attr.NestedAttributes)
		}

		if attr.GenAIValidatorType == "" {
			continue
		}

		slog.Info("Generating validator with OPENAI gpt4o", "validator", attr.GenAIValidatorType)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       openai.GPT4o,
				MaxTokens:   4095,
				TopP:        0.3,
				Temperature: 0.1,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: systemPrompt,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: fmt.Sprintf("type %s struct{} %s", attr.GenAIValidatorType, attr.Description),
					},
				},
			},
		)
		if err != nil {
			slog.Error(fmt.Sprintf("OPENAI error: %s", err), "validator", attr.GenAIValidatorType)
		}

		contents += "\n\n" + cleanResponse(resp.Choices[0].Message.Content) + "\n"

	}
	return contents
}

func cleanResponse(response string) string {
	// Define a regular expression pattern to match the triple backticks and optional language name
	re := regexp.MustCompile("(?s)```[a-zA-Z]*\n(.*)\n```")

	// Find the match and extract the inner content
	matches := re.FindStringSubmatch(response)
	if len(matches) > 1 {
		return matches[1]
	}
	return response
}
