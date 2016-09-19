// Copyright © 2010-12 Qtrac Ltd.
// 
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "fmt"
    "log"
    "math"
    "net/http"
    "sort"
    "strconv"
    "strings"
)

const (
    pageTop    = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>Statistics</title>
<body><h3>Statistics</h3>
<p>Computes basic statistics for a given list of numbers</p>`
    form       = `<form action="/" method="POST">
<label for="numbers">Numbers (comma or space-separated):</label><br />
<input type="text" name="numbers" size="30"><br />
<input type="submit" value="Calculate">
</form>`
    pageBottom = `</body></html>`
    anError    = `<p class="error">%s</p>`
)

type statistics struct {
    numbers []float64
    mean    float64
    median  float64
    stdDev  float64
    mode    float64
}

type pair struct {
    // for the mode...
    Key   string  // save the value of the number as the key
    Value float64 // save the number of times the key is seen as the value
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type pairList []pair
func (p pairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p pairList) Len() int { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func main() {
    http.HandleFunc("/", homePage)
    if err := http.ListenAndServe(":9001", nil); err != nil {
        log.Fatal("failed to start server", err)
    }
}

func homePage(writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm() // Must be called before writing response
    fmt.Fprint(writer, pageTop, form)
    if err != nil {
        fmt.Fprintf(writer, anError, err)
    } else {
        if numbers, message, ok := processRequest(request); ok {
            stats := getStats(numbers)
            fmt.Fprint(writer, formatStats(stats))
        } else if message != "" {
            fmt.Fprintf(writer, anError, message)
        }
    }
    fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([]float64, string, bool) {
    var numbers []float64
    if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
        text := strings.Replace(slice[0], ",", " ", -1)
        for _, field := range strings.Fields(text) {
            if x, err := strconv.ParseFloat(field, 64); err != nil {
                return numbers, "'" + field + "' is invalid", false
            } else {
                numbers = append(numbers, x)    
            }
            
        }
    }
    if len(numbers) == 0 {
        return numbers, "", false // no data first time form is shown
    }
    return numbers, "", true
}

func formatStats(stats statistics) string {
    return fmt.Sprintf(`<table border="1">
<tr><th colspan="2">Results</th></tr>
<tr><td>Numbers</td><td>%v</td></tr>
<tr><td>Count</td><td>%d</td></tr>
<tr><td>Mean</td><td>%f</td></tr>
<tr><td>Median</td><td>%f</td></tr>
<tr><td>StdDev</td><td>%f</td></tr>
<tr><td>Mode</td><td>%f</td></tr>
</table>`, stats.numbers, len(stats.numbers), stats.mean, stats.median, stats.stdDev, stats.mode)
}

func getStats(numbers []float64) (stats statistics) {
    stats.numbers = numbers
    sort.Float64s(stats.numbers)
    stats.mean = sum(numbers) / float64(len(numbers))
    stats.median = median(numbers)
    stats.stdDev = stdDev(numbers, stats.mean)
    stats.mode = mode(numbers)
    return stats
}

func sum(numbers []float64) (total float64) {
    for _, x := range numbers {
        total += x
    }
    return total
}

func median(numbers []float64) float64 {
    middle := len(numbers) / 2
    result := numbers[middle]
    if len(numbers)%2 == 0 {
        result = (result + numbers[middle-1]) / 2
    }
    return result
}

func stdDev(numbers []float64, mean float64) float64 {
    /* std dev calculation:
     * Step 1: Find the mean.
     * Step 2: For each data point, find the square of its distance to the mean.
     * Step 3: Sum the values from Step 2.
     * Step 4: Divide by the number of data points.
     * Step 5: Take the square root.
    */
    var distances []float64
    for _, i := range numbers {
        distance := math.Abs(mean - i)
        square := math.Pow(distance, 2)
        distances = append(distances, square)
    }
    sum := sum(distances)
    x := sum / float64(len(distances))
    result := math.Sqrt(x)
    return result
}

func mode(numbers []float64) float64 {
    /* mode is the most frequently occuring number */
    length := float64(len(numbers))
    if length <= 1 {
        return length 
    }
    frequencyCounts := make(map[string]float64)
    for _, i := range numbers {
        key := strconv.FormatFloat(i, 'f', 1, 64)
        frequencyCounts[key]++
    }
    sortedCounts := sortMapByValue(frequencyCounts)
    mode, err := strconv.ParseFloat(sortedCounts[len(sortedCounts)-1].Key, 64)
    if err != nil {
        log.Fatal(err)
    }
    return mode
}

// A function to turn a map into a pairList, then sort and return it.
func sortMapByValue(m map[string]float64) pairList {
    p := make(pairList, len(m))
    i := 0
    for k, v := range m {
        p[i] = pair{k, v}
        i++
    }
    sort.Sort(p)
    return p
}
