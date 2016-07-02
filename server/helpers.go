package server

import(
  "fmt"
)

func formatByteSize(size int) string {
  sz := float64(size)
  if sz<1024 { return fmt.Sprintf("%db", size) }
  if sz<1024*1024 { return fmt.Sprintf("%.2fk", sz/1024) }
  if sz<1024*1024*1024 { return fmt.Sprintf("%.2fM", sz/1024/1024) }
  if sz<1024*1024*1024 { return fmt.Sprintf("%.2fG", sz/1024/1024/1024) }
  return "LARGE"
} 