package config

// PromptTemplates stores all prompt templates
type PromptTemplates struct {
	Grammar     string `mapstructure:"grammar_template"`
	Vocabulary  string `mapstructure:"vocabulary_template"`
	Translation string `mapstructure:"translation_template"`
	TagExtract  string `mapstructure:"tag_extract_template"`
}

// DefaultPromptTemplates defines default prompt templates
var DefaultPromptTemplates = PromptTemplates{
	Grammar: `You are an English teaching expert. Please provide a detailed grammatical analysis of the following English sentence:

Sentence:
"{{.InputText}}"

Requirements:
- Identify the tense, voice, and grammatical structure used in the sentence
- Explain the grammatical function of each component (subject, predicate, object, etc.)
- Point out and explain any difficult or error-prone points
- Output the explanation in English, suitable for English learners to understand`,

	Vocabulary: `Extract 3~5 core vocabulary words from the following English sentence and output in the following format:

- Word:
- Part of Speech (in English):
- Meaning (in English):
- Example sentence (in English):

Sentence:
"{{.InputText}}"`,

	Translation: `Translate the following English sentence into Chinese and output in two styles:

1. Standard Chinese translation: faithful to the original structure
2. Colloquial Chinese translation: more natural and idiomatic expression

Sentence:
"{{.InputText}}"`,

	TagExtract: `Extract 3-5 representative tags from the following English sentence and analysis results. The tags should be short (1-3 words).
Format: one tag per line, no numbering or other symbols. The result should only contain the tag list, no other explanation.

English sentence: "{{.InputText}}"

Analysis mode: {{.Mode}}

Analysis result:
{{.ResultContent}}`,
}
