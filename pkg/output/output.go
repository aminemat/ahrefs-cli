package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/amine/ahrefs-cli/pkg/client"
)

// Format represents an output format type
type Format string

const (
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
	FormatCSV   Format = "csv"
	FormatTable Format = "table"
)

// Writer handles output formatting and writing
type Writer struct {
	format Format
	writer io.Writer
}

// NewWriter creates a new output writer
func NewWriter(format string, outputFile string) (*Writer, error) {
	var w io.Writer = os.Stdout

	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}
		w = f
	}

	return &Writer{
		format: Format(format),
		writer: w,
	}, nil
}

// WriteSuccess writes a successful response
func (w *Writer) WriteSuccess(data interface{}, meta *client.ResponseMeta) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(data, meta)
	case FormatYAML:
		return w.writeYAML(data, meta)
	case FormatCSV:
		return w.writeCSV(data)
	case FormatTable:
		return w.writeTable(data)
	default:
		return fmt.Errorf("unsupported output format: %s", w.format)
	}
}

// WriteError writes an error response
func (w *Writer) WriteError(err error) error {
	errResp := map[string]interface{}{
		"status": "error",
		"error":  formatError(err),
	}

	enc := json.NewEncoder(w.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(errResp)
}

// writeJSON outputs data as JSON
func (w *Writer) writeJSON(data interface{}, meta *client.ResponseMeta) error {
	response := map[string]interface{}{
		"status": "success",
		"data":   data,
	}

	if meta != nil {
		response["meta"] = map[string]interface{}{
			"response_time_ms": meta.ResponseTimeMS,
		}
		if meta.UnitsConsumed > 0 {
			response["meta"].(map[string]interface{})["units_consumed"] = meta.UnitsConsumed
		}
		if meta.RateLimitRemaining > 0 {
			response["meta"].(map[string]interface{})["rate_limit_remaining"] = meta.RateLimitRemaining
		}
	}

	enc := json.NewEncoder(w.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(response)
}

// writeYAML outputs data as YAML (simple implementation)
func (w *Writer) writeYAML(data interface{}, meta *client.ResponseMeta) error {
	// Simple YAML implementation without external deps
	fmt.Fprintln(w.writer, "status: success")
	fmt.Fprintln(w.writer, "data:")
	return w.writeYAMLValue(data, 1)
}

func (w *Writer) writeYAMLValue(v interface{}, indent int) error {
	prefix := strings.Repeat("  ", indent)

	val := reflect.ValueOf(v)
	if !val.IsValid() {
		fmt.Fprintf(w.writer, "%snil\n", prefix)
		return nil
	}

	switch val.Kind() {
	case reflect.Map:
		for _, key := range val.MapKeys() {
			fmt.Fprintf(w.writer, "%s%v:\n", prefix, key.Interface())
			if err := w.writeYAMLValue(val.MapIndex(key).Interface(), indent+1); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			fmt.Fprintf(w.writer, "%s-\n", prefix)
			if err := w.writeYAMLValue(val.Index(i).Interface(), indent+1); err != nil {
				return err
			}
		}
	case reflect.Struct:
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			if field.IsExported() {
				fmt.Fprintf(w.writer, "%s%s:\n", prefix, field.Name)
				if err := w.writeYAMLValue(val.Field(i).Interface(), indent+1); err != nil {
					return err
				}
			}
		}
	default:
		fmt.Fprintf(w.writer, "%s%v\n", prefix, v)
	}

	return nil
}

// writeCSV outputs data as CSV
func (w *Writer) writeCSV(data interface{}) error {
	csvWriter := csv.NewWriter(w.writer)
	defer csvWriter.Flush()

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Map {
		// If data is a map, try to extract an array/slice field
		for _, key := range val.MapKeys() {
			fieldVal := val.MapIndex(key)
			if fieldVal.Kind() == reflect.Slice || fieldVal.Kind() == reflect.Array {
				val = fieldVal
				break
			}
		}
	}

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return fmt.Errorf("CSV format requires array/slice data")
	}

	if val.Len() == 0 {
		return nil
	}

	// Get headers from first element
	first := val.Index(0)
	headers := extractHeaders(first)
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	// Write rows
	for i := 0; i < val.Len(); i++ {
		row := extractRow(val.Index(i), headers)
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// writeTable outputs data as a formatted table
func (w *Writer) writeTable(data interface{}) error {
	tw := tabwriter.NewWriter(w.writer, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Map {
		// If data is a map, try to extract an array/slice field
		for _, key := range val.MapKeys() {
			fieldVal := val.MapIndex(key)
			if fieldVal.Kind() == reflect.Slice || fieldVal.Kind() == reflect.Array {
				val = fieldVal
				break
			}
		}
	}

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		// Single object - print as key-value pairs
		return w.writeTableObject(tw, data)
	}

	if val.Len() == 0 {
		fmt.Fprintln(tw, "(no results)")
		return nil
	}

	// Get headers
	headers := extractHeaders(val.Index(0))
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	fmt.Fprintln(tw, strings.Repeat("-", len(headers)*10))

	// Write rows
	for i := 0; i < val.Len(); i++ {
		row := extractRow(val.Index(i), headers)
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}

	return nil
}

// writeTableObject writes a single object as a table
func (w *Writer) writeTableObject(tw *tabwriter.Writer, data interface{}) error {
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Map {
		for _, key := range val.MapKeys() {
			fmt.Fprintf(tw, "%v:\t%v\n", key.Interface(), val.MapIndex(key).Interface())
		}
		return nil
	}

	if val.Kind() == reflect.Struct {
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			if field.IsExported() {
				fmt.Fprintf(tw, "%s:\t%v\n", field.Name, val.Field(i).Interface())
			}
		}
		return nil
	}

	fmt.Fprintf(tw, "Value:\t%v\n", data)
	return nil
}

// extractHeaders extracts field names from a value
func extractHeaders(v reflect.Value) []string {
	var headers []string

	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			headers = append(headers, fmt.Sprintf("%v", key.Interface()))
		}
		return headers
	}

	if v.Kind() == reflect.Struct {
		typ := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := typ.Field(i)
			if field.IsExported() {
				// Use JSON tag if available
				jsonTag := field.Tag.Get("json")
				if jsonTag != "" && jsonTag != "-" {
					name := strings.Split(jsonTag, ",")[0]
					headers = append(headers, name)
				} else {
					headers = append(headers, field.Name)
				}
			}
		}
	}

	return headers
}

// extractRow extracts values from a row based on headers
func extractRow(v reflect.Value, headers []string) []string {
	row := make([]string, len(headers))

	if v.Kind() == reflect.Map {
		for i, header := range headers {
			for _, key := range v.MapKeys() {
				if fmt.Sprintf("%v", key.Interface()) == header {
					row[i] = fmt.Sprintf("%v", v.MapIndex(key).Interface())
					break
				}
			}
		}
		return row
	}

	if v.Kind() == reflect.Struct {
		typ := v.Type()
		for i, header := range headers {
			for j := 0; j < v.NumField(); j++ {
				field := typ.Field(j)
				jsonTag := field.Tag.Get("json")
				fieldName := field.Name
				if jsonTag != "" && jsonTag != "-" {
					fieldName = strings.Split(jsonTag, ",")[0]
				}
				if fieldName == header {
					row[i] = fmt.Sprintf("%v", v.Field(j).Interface())
					break
				}
			}
		}
	}

	return row
}

// formatError formats an error for output
func formatError(err error) map[string]interface{} {
	errMap := map[string]interface{}{
		"message": err.Error(),
	}

	// Check if it's an API error
	if apiErr, ok := err.(*client.APIError); ok {
		errMap["code"] = apiErr.Code
		errMap["message"] = apiErr.Message
		if apiErr.Suggestion != "" {
			errMap["suggestion"] = apiErr.Suggestion
		}
		if apiErr.DocsURL != "" {
			errMap["docs_url"] = apiErr.DocsURL
		}
	}

	return errMap
}

// Close closes the writer if it's a file
func (w *Writer) Close() error {
	if f, ok := w.writer.(*os.File); ok && f != os.Stdout {
		return f.Close()
	}
	return nil
}
