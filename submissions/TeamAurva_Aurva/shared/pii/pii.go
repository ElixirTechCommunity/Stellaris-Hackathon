package pii

import (
	"regexp"
	"strconv"
	"strings"
)

// PIIType represents the type of PII detected
type PIIType string

const (
	PIITypeAadhaar  PIIType = "aadhaar"
	PIITypePAN      PIIType = "pan"
	PIITypeGSTIN    PIIType = "gstin"
	PIITypePhone    PIIType = "phone"
	PIITypeVoterID  PIIType = "voter_id"
	PIITypeAccount  PIIType = "bank_account"
)

// Detection represents a PII detection result
type Detection struct {
	Type           PIIType
	Value          string  // Masked value for display
	ConfidenceScore float32
	RiskLevel      string
	DPDPSection    string
}

// Classifier is the proprietary PII detection engine
type Classifier struct {
	aadhaarPattern *regexp.Regexp
	panPattern     *regexp.Regexp
	gstinPattern   *regexp.Regexp
	phonePattern   *regexp.Regexp
	voterIDPattern *regexp.Regexp
	accountPattern *regexp.Regexp
}

// NewClassifier creates a new PII classifier instance
func NewClassifier() *Classifier {
	return &Classifier{
		// Aadhaar: 12 digits, first digit is 2-9, optionally with spaces/dashes
		aadhaarPattern: regexp.MustCompile(`\b[2-9]\d{3}[\s-]?\d{4}[\s-]?\d{4}\b`),
		
		// PAN: 5 letters, 4 digits, 1 letter (e.g., ABCDE1234F)
		panPattern: regexp.MustCompile(`\b[A-Z]{3}[PCHABGJLFT][A-Z]\d{4}[A-Z]\b`),
		
		// GSTIN: 15 characters (2 state code + 10 PAN + 1 entity + 1 Z + 1 checksum)
		gstinPattern: regexp.MustCompile(`\b\d{2}[A-Z]{5}\d{4}[A-Z][1-9A-Z]Z[0-9A-Z]\b`),
		
		// Indian phone: 10 digits, optional +91 prefix
		phonePattern: regexp.MustCompile(`(?:\+91[\s-]|91[\s-]|^|[\s,;:"|'])([6-9]\d{9})\b`),
		
		// Voter ID: 3 letters followed by 7 digits (e.g., ABC1234567)
		voterIDPattern: regexp.MustCompile(`\b[A-Z]{3}\d{7}\b`),

		// Bank Account: 9-18 digits, usually found with context
		accountPattern: regexp.MustCompile(`(?i)(?:account|acct|acc|a/c|bank|saving|current)[^0-9\n]{0,20}(\b\d{9,18}\b)`),
	}
}

// Scan analyzes text and returns all PII detections
func (c *Classifier) Scan(text string) []Detection {
	var detections []Detection
	
	// Check for Aadhaar
	if matches := c.aadhaarPattern.FindAllString(text, -1); len(matches) > 0 {
		for _, match := range matches {
			normalized := strings.ReplaceAll(strings.ReplaceAll(match, " ", ""), "-", "")
			if c.validateAadhaar(normalized) {
				detections = append(detections, Detection{
					Type:            PIITypeAadhaar,
					Value:           c.maskAadhaar(normalized),
					ConfidenceScore: 0.95,
					RiskLevel:       "critical",
					DPDPSection:     "Section 8(4) - Sensitive Personal Data",
				})
			}
		}
	}

	// Check for Bank Account (Higher priority if context is present)
	if matches := c.accountPattern.FindAllStringSubmatch(text, -1); len(matches) > 0 {
		for _, match := range matches {
			if len(match) > 1 {
				val := match[1]
				// Deduplicate if already caught as Aadhaar (prefer Bank Account if context matches)
				isDuplicate := false
				for i, d := range detections {
					if d.Type == PIITypeAadhaar && strings.HasSuffix(val, d.Value[len(d.Value)-4:]) {
						detections[i] = Detection{
							Type:            PIITypeAccount,
							Value:           c.maskAccount(val),
							ConfidenceScore: 0.90,
							RiskLevel:       "high",
							DPDPSection:     "Section 8(3) - Financial Data",
						}
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					detections = append(detections, Detection{
						Type:            PIITypeAccount,
						Value:           c.maskAccount(val),
						ConfidenceScore: 0.88,
						RiskLevel:       "high",
						DPDPSection:     "Section 8(3) - Financial Data",
					})
				}
			}
		}
	}
	
	// Check for PAN
	if matches := c.panPattern.FindAllString(text, -1); len(matches) > 0 {
		for _, match := range matches {
			if c.validatePAN(match) {
				detections = append(detections, Detection{
					Type:            PIITypePAN,
					Value:           c.maskPAN(match),
					ConfidenceScore: 0.92,
					RiskLevel:       "high",
					DPDPSection:     "Section 8(3) - Financial Data",
				})
			}
		}
	}
	
	// Check for GSTIN
	if matches := c.gstinPattern.FindAllString(text, -1); len(matches) > 0 {
		for _, match := range matches {
			detections = append(detections, Detection{
				Type:            PIITypeGSTIN,
				Value:           c.maskGSTIN(match),
				ConfidenceScore: 0.88,
				RiskLevel:       "medium",
				DPDPSection:     "Section 7(1) - Business Data",
			})
		}
	}
	
	// Check for Phone
	if matches := c.phonePattern.FindAllStringSubmatch(text, -1); len(matches) > 0 {
		for _, match := range matches {
			if len(match) > 1 {
				val := match[1]
				detections = append(detections, Detection{
					Type:            PIITypePhone,
					Value:           c.maskPhone(val),
					ConfidenceScore: 0.85,
					RiskLevel:       "high",
					DPDPSection:     "Section 8(1) - Contact Information",
				})
			}
		}
	}
	
	// Check for Voter ID
	if matches := c.voterIDPattern.FindAllString(text, -1); len(matches) > 0 {
		for _, match := range matches {
			detections = append(detections, Detection{
				Type:            PIITypeVoterID,
				Value:           c.maskVoterID(match),
				ConfidenceScore: 0.80,
				RiskLevel:       "high",
				DPDPSection:     "Section 8(2) - Government ID",
			})
		}
	}
	
	return detections
}

// validateAadhaar validates Aadhaar using Verhoeff checksum algorithm
func (c *Classifier) validateAadhaar(aadhaar string) bool {
	if len(aadhaar) != 12 {
		return false
	}
	
	// Verhoeff algorithm tables
	d := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{1, 2, 3, 4, 0, 6, 7, 8, 9, 5},
		{2, 3, 4, 0, 1, 7, 8, 9, 5, 6},
		{3, 4, 0, 1, 2, 8, 9, 5, 6, 7},
		{4, 0, 1, 2, 3, 9, 5, 6, 7, 8},
		{5, 9, 8, 7, 6, 0, 4, 3, 2, 1},
		{6, 5, 9, 8, 7, 1, 0, 4, 3, 2},
		{7, 6, 5, 9, 8, 2, 1, 0, 4, 3},
		{8, 7, 6, 5, 9, 3, 2, 1, 0, 4},
		{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	}
	
	p := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{1, 5, 7, 6, 2, 8, 3, 0, 9, 4},
		{5, 8, 0, 3, 7, 9, 6, 1, 4, 2},
		{8, 9, 1, 6, 0, 4, 3, 5, 2, 7},
		{9, 4, 5, 3, 1, 2, 6, 8, 7, 0},
		{4, 2, 8, 6, 5, 7, 3, 9, 0, 1},
		{2, 7, 9, 3, 8, 0, 6, 4, 1, 5},
		{7, 0, 4, 6, 9, 1, 3, 2, 5, 8},
	}
	
	c_val := 0
	for i := 0; i < len(aadhaar); i++ {
		digit, _ := strconv.Atoi(string(aadhaar[len(aadhaar)-1-i]))
		c_val = d[c_val][p[(i)%8][digit]]
	}
	
	return c_val == 0
}

// validatePAN validates PAN card format - 4th character validation
func (c *Classifier) validatePAN(pan string) bool {
	if len(pan) != 10 {
		return false
	}
	
	// 4th character indicates entity type
	fourthChar := pan[3]
	validChars := map[byte]bool{
		'P': true, // Individual
		'C': true, // Company
		'H': true, // HUF
		'A': true, // AOP
		'B': true, // BOI
		'G': true, // Government
		'J': true, // Artificial Juridical Person
		'L': true, // Local Authority
		'F': true, // Firm
		'T': true, // Trust
	}
	
	return validChars[fourthChar]
}

// Masking functions
func (c *Classifier) maskAadhaar(aadhaar string) string {
	if len(aadhaar) != 12 {
		return "XXXX-XXXX-XXXX"
	}
	return "XXXX-XXXX-" + aadhaar[8:]
}

func (c *Classifier) maskPAN(pan string) string {
	if len(pan) != 10 {
		return "XXXXX1234X"
	}
	return "XXXXX" + pan[5:9] + "X"
}

func (c *Classifier) maskGSTIN(gstin string) string {
	if len(gstin) != 15 {
		return "XXXXXXXXXXXX..."
	}
	return "XX" + gstin[2:7] + "XXXXXXXX"
}

func (c *Classifier) maskPhone(phone string) string {
	digits := regexp.MustCompile(`\d`).FindAllString(phone, -1)
	if len(digits) >= 10 {
		lastDigits := strings.Join(digits[len(digits)-10:], "")
		return "XXXXXX" + lastDigits[6:]
	}
	return "XXXXXXX" + phone[len(phone)-3:]
}

func (c *Classifier) maskVoterID(voterID string) string {
	if len(voterID) != 10 {
		return "XXXXX..."
	}
	return "XXX" + voterID[3:7] + "XXX"
}

func (c *Classifier) maskAccount(account string) string {
	if len(account) < 4 {
		return "****"
	}
	return "XXXXXX" + account[len(account)-4:]
}
