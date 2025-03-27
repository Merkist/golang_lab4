package main

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}

type result struct {
	Key   string
	Value float64
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/part1", func(c *gin.Context) {
		c.HTML(http.StatusOK, "part1.html", nil)
	})

	r.GET("/part2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "part2.html", nil)
	})

	r.GET("/part3", func(c *gin.Context) {
		c.HTML(http.StatusOK, "part3.html", nil)
	})

	r.POST("/calculate_part1", func(c *gin.Context) {
		values := []string{"Voltage", "Current", "FicTime", "Load", "Time"}
		inputs := make(map[string]float64)
		for _, v := range values {
			val, err := strconv.ParseFloat(c.PostForm(v), 64)
			if err != nil {
				c.HTML(http.StatusBadRequest, "part1.html", gin.H{"error": "Invalid input for " + v})
				return
			}
			inputs[v] = val
		}

		Im := (inputs["Load"] / 2) / (math.Sqrt(3) * inputs["Voltage"])
		Impa := Im * 2

		var jEc float64
		if inputs["Time"] >= 1000 && inputs["Time"] <= 3000 {
			jEc = 1.6
		} else if inputs["Time"] > 3000 && inputs["Time"] <= 5000 {
			jEc = 1.4
		} else if inputs["Time"] > 5000 {
			jEc = 1.2
		}

		cT := 92.0

		sEconom := Im / jEc
		sMin := (inputs["Current"] * 1000 * math.Sqrt(inputs["FicTime"])) / cT

		c.HTML(http.StatusOK, "part1.html", gin.H{
			"CurrentNormal": roundToTwoDecimalPlaces(Im),
			"CurrentEmerg":  roundToTwoDecimalPlaces(Impa),
			"SectionEconom": roundToTwoDecimalPlaces(sEconom),
			"SectionMin":    roundToTwoDecimalPlaces(sMin),
		})
	})

	r.POST("/calculate_part2", func(c *gin.Context) {
		values := []string{"Voltage2", "Power"}
		inputs := make(map[string]float64)
		for _, v := range values {
			val, err := strconv.ParseFloat(c.PostForm(v), 64)
			if err != nil {
				c.HTML(http.StatusBadRequest, "part2.html", gin.H{"error": "Invalid input for " + v})
				return
			}
			inputs[v] = val
		}

		sNomT := 6.3

		xC := math.Pow(inputs["Voltage2"], 2) / inputs["Power"]
		xT := (inputs["Voltage2"] / 100) * (math.Pow(inputs["Voltage2"], 2) / sNomT)
		x := xC + xT
		iP0 := inputs["Voltage2"] / (math.Sqrt(3) * x)

		c.HTML(http.StatusOK, "part2.html", gin.H{
			"resultResistance": roundToTwoDecimalPlaces(x),
			"resultCurrent":    roundToTwoDecimalPlaces(iP0),
		})
	})

	r.POST("/calculate_part3", func(c *gin.Context) {
		values := []string{"ResistanceNormR", "ResistanceNormX", "ResistanceMinR", "ResistanceMinX"}
		inputs := make(map[string]float64)
		for _, v := range values {
			val, err := strconv.ParseFloat(c.PostForm(v), 64)
			if err != nil {
				c.HTML(http.StatusBadRequest, "part3.html", gin.H{"error": "Invalid input for " + v})
				return
			}
			inputs[v] = val
		}

		Ukmax := 11.1
		Uvn := 115.0
		SnomT := 6.3

		xT := (Ukmax * Uvn * Uvn) / (100 * SnomT)
		xH := inputs["ResistanceNormX"] + xT
		zH := math.Sqrt(math.Pow(inputs["ResistanceNormR"], 2) + math.Pow(xH, 2))
		xHMin := inputs["ResistanceMinX"] + xT
		zHMin := math.Sqrt(math.Pow(inputs["ResistanceMinR"], 2) + math.Pow(xHMin, 2))

		iH3 := (Uvn * 1000) / (math.Sqrt(3) * zH)
		iH2 := iH3 * (math.Sqrt(3) / 2)
		iH3Min := (Uvn * 1000) / (math.Sqrt(3) * zHMin)
		iH2Min := iH3Min * (math.Sqrt(3) / 2)

		Unn := 11.0
		k := math.Pow(Unn, 2) / math.Pow(Uvn, 2)

		rHN := inputs["ResistanceNormR"] * k
		xHN := xH * k
		zHN := math.Sqrt(math.Pow(rHN, 2) + math.Pow(xHN, 2))
		rHNMin := inputs["ResistanceMinR"] * k
		xHNMin := xHMin * k
		zHNMin := math.Sqrt(math.Pow(rHNMin, 2) + math.Pow(xHNMin, 2))

		iHN3 := (Unn * 1000) / (math.Sqrt(3) * zHN)
		iHN2 := iHN3 * (math.Sqrt(3) / 2)
		iHN3Min := (Unn * 1000) / (math.Sqrt(3) * zHNMin)
		iHN2Min := iHN3Min * (math.Sqrt(3) / 2)

		results_1 := []result{
			{"Струм трифазного K3 в нормальному режимі:", roundToTwoDecimalPlaces(iH3)},
			{"Струм двофазного K3 в нормальному режимі:", roundToTwoDecimalPlaces(iH2)},
			{"Струм трифазного K3 в мінімальному режимі:", roundToTwoDecimalPlaces(iH3Min)},
			{"Струм двофазного K3 в мінімальному режимі:", roundToTwoDecimalPlaces(iH2Min)},
		}
		results_2 := []result{
			{"Струм трифазного K3 в нормальному режимі:", roundToTwoDecimalPlaces(iHN3)},
			{"Струм двофазного K3 в нормальному режимі:", roundToTwoDecimalPlaces(iHN2)},
			{"Струм трифазного K3 в мінімальному режимі:", roundToTwoDecimalPlaces(iHN3Min)},
			{"Струм двофазного K3 в мінімальному режимі:", roundToTwoDecimalPlaces(iHN2Min)},
		}

		c.HTML(http.StatusOK, "part3.html", gin.H{
			"results_1": results_1,
			"results_2": results_2,
		})
	})

	r.Run(":8080")
}
