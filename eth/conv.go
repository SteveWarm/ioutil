package eth

import (
    "fmt"
    "strconv"
    "strings"
)

// ipv4格式的地址字符串转uint32
func IPv4StrToUint32(ip string) uint32 {
    bits := strings.Split(ip, ".")
    b0, _ := strconv.Atoi(bits[0])
    b1, _ := strconv.Atoi(bits[1])
    b2, _ := strconv.Atoi(bits[2])
    b3, _ := strconv.Atoi(bits[3])
    var sum uint32
    sum += uint32(uint8(b0)) << 24
    sum += uint32(uint8(b1)) << 16
    sum += uint32(uint8(b2)) << 8
    sum += uint32(uint8(b3))
    return sum
}

// ipv4 uint32转字符串
func Uint32ToIPv4Str(intip uint32) string {
    ip1 := intip >> 24
    ip2 := (intip >> 16) & 0x000000FF
    ip3 := (intip >> 8) & 0x000000FF
    ip4 := intip & 0x000000FF
    return fmt.Sprintf("%d.%d.%d.%d", ip1, ip2, ip3, ip4)
}
