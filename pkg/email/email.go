package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/llamacto/llama-gin-kit/config"
	"github.com/llamacto/llama-gin-kit/pkg/logger"
)

var (
	cfg *config.Config
)

// Init 初始化邮件服务
func Init(c *config.Config) {
	cfg = c
}

type EmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

type EmailResponse struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Created string `json:"created"`
	Error   string `json:"error"`
}

// SendEmail 发送邮件
func SendEmail(to []string, subject, htmlContent string) error {
	if cfg == nil {
		return fmt.Errorf("email service not initialized")
	}

	logger.Info("Preparing to send email",
		fmt.Sprintf("from: %s", cfg.Email.From),
		fmt.Sprintf("to: %v", to),
		fmt.Sprintf("subject: %s", subject),
	)

	reqBody := EmailRequest{
		From:    cfg.Email.From,
		To:      to,
		Subject: subject,
		Html:    htmlContent,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		logger.Error("Failed to serialize request", err)
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	logger.Info("Request data", string(jsonData))

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Failed to create request", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Email.ResendAPIKey)
	logger.Info("Using API Key", cfg.Email.ResendAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to send request", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	logger.Info("Received response", string(body))

	if resp.StatusCode == http.StatusForbidden {
		var resendError struct {
			Name       string `json:"name"`
			Message    string `json:"message"`
			StatusCode int    `json:"statusCode"`
		}
		if err := json.Unmarshal(body, &resendError); err != nil {
			logger.Error("Failed to parse error response", err)
			return fmt.Errorf("failed to unmarshal error response: %w", err)
		}
		logger.Error("Resend API error", fmt.Errorf("%s: %s (status %d)",
			resendError.Name,
			resendError.Message,
			resendError.StatusCode,
		))
		if resendError.Name == "validation_error" && strings.Contains(resendError.Message, "domain is not verified") {
			return fmt.Errorf("recipient domain not verified, please contact admin to add domain verification")
		}
		return fmt.Errorf("Resend API error: %s", resendError.Message)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.Error("Email sending failed", fmt.Errorf("status code %d, response: %s", resp.StatusCode, string(body)))
		return fmt.Errorf("failed to send email: status code %d, response: %s", resp.StatusCode, string(body))
	}

	var emailResp EmailResponse
	if err := json.Unmarshal(body, &emailResp); err != nil {
		logger.Error("Failed to parse response", err)
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if emailResp.Error != "" {
		logger.Error("Email service error", fmt.Errorf(emailResp.Error))
		return fmt.Errorf("email service error: %s", emailResp.Error)
	}

	logger.Info("Email sent successfully", fmt.Sprintf("ID: %s", emailResp.ID))
	return nil
}

// SendPasswordResetEmail sends a password reset notification email
func SendPasswordResetEmail(to string, newPassword string) error {
	subject := "Password Reset Notification"
	htmlContent := fmt.Sprintf(`
		<h2>Password Reset Notification</h2>
		<p>Your password has been reset. The new temporary password is:</p>
		<p style="font-size: 18px; font-weight: bold; color: #333;">%s</p>
		<p>Please use this temporary password to log in and change it to your own password immediately.</p>
		<p>If this was not your action, please contact the administrator immediately.</p>
	`, newPassword)

	return SendEmail([]string{to}, subject, htmlContent)
}

// SendWelcomeEmail sends a welcome email
func SendWelcomeEmail(to string, username string) error {
	subject := "Welcome to Llama Gin Kit"
	htmlContent := fmt.Sprintf(`
		<h2>Welcome to Llama Gin Kit</h2>
		<p>Dear %s,</p>
		<p>Thank you for registering as our user!</p>
		<p>If you have any questions, please feel free to contact our support team.</p>
	`, username)

	return SendEmail([]string{to}, subject, htmlContent)
}
