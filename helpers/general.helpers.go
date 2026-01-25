package helpers

import (
	"crypto/rand"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	// BEGIN __INCLUDE_TEMPLATE__
	"fmt"
	"math"
	"math/big"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/leekchan/accounting"
	// END __INCLUDE_TEMPLATE__
)

// LoadEnv ...
func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Warn("Error loading environtment from file .env")
		log.Info("Loading environment from local machine ...")
		return
	}
	log.Info("Loading environment from .env file ...")
}

// GetEnv ..
func GetEnv(key string) string {

	// load .env file

	// if err != nil {
	// 	os.Setenv(key, defaultString)
	// }

	onceLoadEnv.Do(func() {
		LoadEnv()
	})

	value, present := os.LookupEnv(key)
	if !present {
		log.Fatalf("ENV %v not found", key)
	}
	// value := os.Getenv(key)
	// if value == "" {
	// 	log.Fatalf("ENV %v not found", key)
	// }

	return value
	// viper.SetConfigType("env")
	// viper.AddConfigPath(".")
	// viper.SetConfigName(".env")
	// viper.AutomaticEnv()
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	log.Fatalf("Error while reading config file %s", err)
	// }
	// value, ok := viper.Get(key).(string)
	// if !ok {
	// 	log.Fatalf("Invalid type assertion %s", key)
	// }

	// return value
}

// BEGIN __INCLUDE_TEMPLATE__

// PrintHeader ...
func PrintHeader() {
	pc, _, _, _ := runtime.Caller(1)
	fmt.Printf("<======> %s <======>", runtime.FuncForPC(pc).Name())
	fmt.Println()
}

// DiffTime ...
func DiffTime(a, b time.Time) string {
	var age string
	// if a.Location() != b.Location() {
	// 	b = b.In(a.Location())
	// }
	// loc, _ := time.LoadLocation("Asia/Jakarta")
	// a = a.In(loc)
	// b = b.In(loc)
	// b = b.UTC()

	locTime := b.In(a.Location())
	_, zoneOffset := locTime.Zone()

	b = locTime.Add(-time.Duration(zoneOffset) * time.Second)

	fmt.Println("a => ", a)
	fmt.Println("b => ", b)

	if a.After(b) {
		a, b = b, a
	}

	diff := math.RoundToEven(b.Sub(a).Hours())
	day := int64(diff / 24)
	hour := int(diff) % 24

	if day == 0 {
		age = fmt.Sprintf("%v hours", hour)
	} else {
		age = fmt.Sprintf("%v days %v hours", day, hour)
	}

	return age
	// return b.Local().Sub(a).String()
}

// DiffTimes ...
func DiffTimes(a, b time.Time) (string, string, string) {
	nano := b.Sub(a)
	return a.Format("2006-01-02 15:04"), b.Format("2006-01-02 15:04"), nano.String()
}

// ConvertStringToDate ...
func ConvertStringToDate(s, layoutISO string) (time.Time, error) {
	t, err := time.Parse(layoutISO, s)
	if err != nil {
		fmt.Println("Error Convert Date => ", err)
		return t, err
	}
	return t, nil
}

// CheckArray ...
func CheckArray(data string, arrayCheck []string) bool {
	for _, v := range arrayCheck {
		if v == data {
			return true
		}
	}
	return false
}

// CheckDir ...
func CheckDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Println("====== Error os.Stat ======")
		fmt.Printf("====== %v ======", err)
	}

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, os.ModePerm)
		if errDir != nil {
			return errDir
		}

	}
	return nil
}

// GetCurrency ..
func GetCurrency(current float64) string {
	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2, Thousand: ".", Decimal: ","}
	data := big.NewFloat(current)
	return ac.FormatMoneyBigFloat(data)
}

// ArrayStringToArrayInt converts a string array to an array of ints.
func ArrayStringToArrayInt(data []string) ([]int, error) {
	var result []int
	for _, v := range data {
		i, err := strconv.Atoi(v)
		if err != nil {
			return result, err
		}
		result = append(result, i)
	}
	return result, nil
}

// ErrorValidator ..
func ErrorValidator(err error) string {
	var errors, message string
	for _, err := range err.(validator.ValidationErrors) {
		// fmt.Println(err.Namespace())
		// fmt.Println(err.StructNamespace())
		// fmt.Println(err.StructField())
		// fmt.Println(err.Tag())
		// fmt.Println(err.ActualTag())
		// fmt.Println(err.Kind())
		// fmt.Println(err.Type())
		// fmt.Println(err.Value())
		// fmt.Println(err.Param())
		// fmt.Println()
		switch err.Tag() {
		case "required":
			message = "Tidak Boleh Kosong"
		case "numeric":
			message = "Harus Berupa Angka"
		case "ne":
			message = "Karakter Tidak Diperbolehkan"
		case "date-time":
			message = "Format Tanggal Salah, Silahkan Ikuti Format 2006-01-02 15:04"
		}

		errors += fmt.Sprintf("%s %s, ", err.StructField(), message)
	}
	return errors
}

// TrimZero trims a zero - length string
func TrimZero(data string) string {
	return strings.TrimRight(data, "0")
}

// SimiliarTo translates a slice of strings into a string.
func SimiliarTo(data []string) string {
	return "(" + strings.Join(data, "|") + ")%"
}

// UniqueArray ...
func UniqueArray(data []string) []string {
	m := make(map[string]bool)
	for _, v := range data {
		if _, ok := m[v]; !ok {
			m[v] = true
		}
	}
	var result []string
	for k := range m {
		result = append(result, k)
	}
	return result
}

// BeginningOfMonth ...
func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

// EndOfMonth ...
func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}

// HandleNan ...
func HandleNan(data float64) float64 {
	if math.IsNaN(data) {
		data = 0
	}

	return data
}

// RoundFloat ..
func RoundFloat(input float64, precision int) float64 {
	temp := float64(1)
	for i := 0; i < precision; i++ {
		temp *= 10
	}
	return math.Round(input*temp) / temp
}

func JoinInts(ints []int, sep string) string {
	strs := make([]string, len(ints))
	for i, val := range ints {
		strs[i] = strconv.Itoa(val)
	}
	return strings.Join(strs, sep)
}

// MonthsElapsed menghitung jumlah bulan penuh yang telah berlalu
// antara waktu t1 (awal/createdAt) dan t2 (akhir/sekarang).
func MonthsElapsed(t1, t2 time.Time) int {
    // 1. Pastikan t2 (waktu akhir) lebih lambat dari t1 (waktu awal)
	if t1.After(t2) {
		t1, t2 = t2, t1 // Swap jika t1 lebih baru dari t2
	}

    // 2. Ekstrak komponen Tahun dan Bulan
	year1, month1, _ := t1.Date()
	year2, month2, _ := t2.Date()
	
    // 3. Hitung selisih Tahun dalam satuan bulan
    // Contoh: 2025 - 2024 = 1 tahun * 12 = 12 bulan
	diffYears := year2 - year1
	totalMonths := diffYears * 12
	
    // 4. Tambahkan/kurangi selisih Bulan
	totalMonths += int(month2) - int(month1)

    // 5. Penyesuaian Harian (Mengecek Bulan Penuh)
    // Jika hari di t2 (waktu akhir) lebih kecil dari hari di t1 (waktu awal), 
    // maka bulan terakhir belum terhitung penuh. Kurangi 1 bulan.
    // Contoh: 31 Jan (t1) ke 15 Feb (t2). Selisih total 1 bulan. Namun,
    // hari (15) < hari (31), jadi bulan Februari belum penuh. totalMonths - 1.
    // Contoh: 31 Jan (t1) ke 31 Mar (t2). Selisih total 2 bulan. Hari (31) >= Hari (31).
	if t2.Day() < t1.Day() {
		totalMonths--
	}
	
	return totalMonths
}


const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString menghasilkan string acak alfanumerik sepanjang n karakter
func RandomString(n int) string {
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// fallback kalau error (sangat jarang terjadi)
			result[i] = charset[0]
			continue
		}
		result[i] = charset[num.Int64()]
	}
	return string(result)
}




// END __INCLUDE_TEMPLATE__
