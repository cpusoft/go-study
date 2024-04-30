//go:build prod

package main

func init() {
	configArr = append(configArr, "mysql prod")
}
