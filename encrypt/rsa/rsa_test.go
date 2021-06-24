package rsa

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRsa(t *testing.T) {
	priKey, pubKey, _ := GenRsaKey(2048)
	text := "MIIEogIBAAKCAQEApXeHUE6EbVojgcXRgfJdQqfYHCLY6Rp54xf6hJ2h7H0wNSDiQvAWjeQoLRXV3abdL/rP/6hQKG27AHJO78XZVqssW+8WMbKw85Q41Z3bLK1ptQBGAbAC8Ggg5lYdJDFjcvWbB0T+EHRnhjNCWjwI418GtS6JCBewNAZs618cpmLY7U53avCOu39W4C6OPILFfGIo2jXDaw2jAFHgMvVNBohwC0GaJ0XrmW3vcq9qZi1Nm8Ak/uMGTVb5I9vgNcrtfYujIjbpu95ww064/3A7c2JJ/Wb2+5hA4DS5jt6Bak5IIrA2NRBIJJRbtcpMelVJAhSBlNqZcd5gffHPssWQUQIDAQABAoIBAHd3Hj4wAlLFcr9eaM5Og2D9t/1Qd5WNZRU+GcSn3uHU3Ppu3I2BtHcLlKpAiqG2uRnMF2K5Te5yk0JWRYG+MhXuDl/t9fM51aJ3kLVCfJz8M0bYhLLxNp4GQEKtR+r6sZetlhmgiWKt+JSe67gkgjJPSJOFnrA2EiTtwrQJfrDtogJqPqVU4jbJPmacIzSbjNFMfoqhk9Cr15l52AK77U9apWu3aVwufyfFK/0dA9niZmR4lj/mGZ/hWpKF847CakpFvIu+khTUlgvMqXwuuGrQdea8rmJoq7A7MYIxKxFMY7cKWpnhgQPRaVoMJoK4Q5lxFpucy7QyPrcCgp90vMUCgYEAyxH58HJTmBMeQC+gMBMKtsSLM3TQQGErTnZYrHPT/besDPzr5PFwDu1Zt8l+Yb0puKwN4fx6MaIUGbtaA+YQooayfFgiJa2P7o5WSYIVW1y8bZSrrUawGcMIUSl8V1kaJoQvSe1oyKFRaT8kkHOTFTX2walWsoWnY9prNcMCrcsCgYEA0Jhw6NBJhGjrThTP0qXC80GUSyNbiGeALL6NzgCgtNCd5MFtk+gVoa/9Ji26YsDcfwOwBJqrkG3DI1dBn5qUYIarqv+rcFAU6pz2uAmW3JiOQd4jQ0MyqcSsEM5qa0fYamFUKlReMx4IO8KNKyp1BFsYq0wm83mFkRoXeE/GttMCgYBJmAczq8s91tfkvR3Zrlz4pbwo9tGuM0jlk6BJR2Txk0oIHvVCsHlC/6O/Jofl1g8zvS7+0mhaannMZYiW1x76N8ShqbMeYotCElWVKE6jILWtJO8eyfpyK6ts9pL4ePMwOEGHEkIiS8xcTyTqMOiCDF+UCdHAuw1R88tc3YKwBwKBgHUZzvDz1QGvQMGRr2WKxtl2rEBONhlqOStlQggulAlNwAXmjJRmypX9TTj8nNDJgj8Pm+XJypyG8fBKEL3/smJJ199kLiMb4dIfkeWZBIcMYXgas2MUO0HQ9eNtbZKSP6zgvLYSrNs3ddnOix97czuhxESNuKQgSVo+8oQJDP4fAoGATr7dQBaNqrECjKs3CQpfc7FcgmUjqUeeNcufKkOi9cYeZYtHk223NTLZ7FhEEuuBhr0wu2gNqJojZPtt7f2iafoKSyguCENFQCSsXJeCXE8Vgjms198/i393NpBPta15/4PzXmgPzeyfloW8UhrgYYBrgNTExfbrNmXmEutUZaY="
	for i := 0; i < 4096; i++ {
		EncryptStr, err := Encrypt([]byte(text), pubKey)
		if err != nil {
			fmt.Printf("%v \n", err)
		}
		Decrypt, err := Decrypt(EncryptStr, priKey)
		if err != nil {
			fmt.Printf("%v \n", err)
		}

		if string(Decrypt) == text {
			fmt.Printf("rsa parse success!, i=%d, txt=%s", i, text)
		} else {
			fmt.Println("rsa parse Fail!")
		}
		text = text + "x"
		time.Sleep(time.Duration(20000) * time.Nanosecond)
	}
}

func TestDecrypt(t *testing.T) {
	priKey := "MIIEogIBAAKCAQEApXeHUE6EbVojgcXRgfJdQqfYHCLY6Rp54xf6hJ2h7H0wNSDiQvAWjeQoLRXV3abdL/rP/6hQKG27AHJO78XZVqssW+8WMbKw85Q41Z3bLK1ptQBGAbAC8Ggg5lYdJDFjcvWbB0T+EHRnhjNCWjwI418GtS6JCBewNAZs618cpmLY7U53avCOu39W4C6OPILFfGIo2jXDaw2jAFHgMvVNBohwC0GaJ0XrmW3vcq9qZi1Nm8Ak/uMGTVb5I9vgNcrtfYujIjbpu95ww064/3A7c2JJ/Wb2+5hA4DS5jt6Bak5IIrA2NRBIJJRbtcpMelVJAhSBlNqZcd5gffHPssWQUQIDAQABAoIBAHd3Hj4wAlLFcr9eaM5Og2D9t/1Qd5WNZRU+GcSn3uHU3Ppu3I2BtHcLlKpAiqG2uRnMF2K5Te5yk0JWRYG+MhXuDl/t9fM51aJ3kLVCfJz8M0bYhLLxNp4GQEKtR+r6sZetlhmgiWKt+JSe67gkgjJPSJOFnrA2EiTtwrQJfrDtogJqPqVU4jbJPmacIzSbjNFMfoqhk9Cr15l52AK77U9apWu3aVwufyfFK/0dA9niZmR4lj/mGZ/hWpKF847CakpFvIu+khTUlgvMqXwuuGrQdea8rmJoq7A7MYIxKxFMY7cKWpnhgQPRaVoMJoK4Q5lxFpucy7QyPrcCgp90vMUCgYEAyxH58HJTmBMeQC+gMBMKtsSLM3TQQGErTnZYrHPT/besDPzr5PFwDu1Zt8l+Yb0puKwN4fx6MaIUGbtaA+YQooayfFgiJa2P7o5WSYIVW1y8bZSrrUawGcMIUSl8V1kaJoQvSe1oyKFRaT8kkHOTFTX2walWsoWnY9prNcMCrcsCgYEA0Jhw6NBJhGjrThTP0qXC80GUSyNbiGeALL6NzgCgtNCd5MFtk+gVoa/9Ji26YsDcfwOwBJqrkG3DI1dBn5qUYIarqv+rcFAU6pz2uAmW3JiOQd4jQ0MyqcSsEM5qa0fYamFUKlReMx4IO8KNKyp1BFsYq0wm83mFkRoXeE/GttMCgYBJmAczq8s91tfkvR3Zrlz4pbwo9tGuM0jlk6BJR2Txk0oIHvVCsHlC/6O/Jofl1g8zvS7+0mhaannMZYiW1x76N8ShqbMeYotCElWVKE6jILWtJO8eyfpyK6ts9pL4ePMwOEGHEkIiS8xcTyTqMOiCDF+UCdHAuw1R88tc3YKwBwKBgHUZzvDz1QGvQMGRr2WKxtl2rEBONhlqOStlQggulAlNwAXmjJRmypX9TTj8nNDJgj8Pm+XJypyG8fBKEL3/smJJ199kLiMb4dIfkeWZBIcMYXgas2MUO0HQ9eNtbZKSP6zgvLYSrNs3ddnOix97czuhxESNuKQgSVo+8oQJDP4fAoGATr7dQBaNqrECjKs3CQpfc7FcgmUjqUeeNcufKkOi9cYeZYtHk223NTLZ7FhEEuuBhr0wu2gNqJojZPtt7f2iafoKSyguCENFQCSsXJeCXE8Vgjms198/i393NpBPta15/4PzXmgPzeyfloW8UhrgYYBrgNTExfbrNmXmEutUZaY="
	//	pubKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApXeHUE6EbVojgcXRgfJdQqfYHCLY6Rp54xf6hJ2h7H0wNSDiQvAWjeQoLRXV3abdL/rP/6hQKG27AHJO78XZVqssW+8WMbKw85Q41Z3bLK1ptQBGAbAC8Ggg5lYdJDFjcvWbB0T+EHRnhjNCWjwI418GtS6JCBewNAZs618cpmLY7U53avCOu39W4C6OPILFfGIo2jXDaw2jAFHgMvVNBohwC0GaJ0XrmW3vcq9qZi1Nm8Ak/uMGTVb5I9vgNcrtfYujIjbpu95ww064/3A7c2JJ/Wb2+5hA4DS5jt6Bak5IIrA2NRBIJJRbtcpMelVJAhSBlNqZcd5gffHPssWQUQIDAQAB"

	text := "我爱中国！hello china! ..11..11333;;"

	priKeyBytes, err := base64.StdEncoding.DecodeString(priKey)
	assert.Nil(t, err)

	decText := "NPywmyoz3TnLEmHtwAXsJD6RRmt1MUlX2JNqgP64cT3Xk33EKdCHDQgBRNRQ2NVEom4IuzyCMPNT1hEPwqua4fUEnA01IOeSR4mDu/uPuN7ULXmTfRrehwdV5Vy1bvFPesEXlOvotGmwVgypowqrBbiA77SAFL/ejdIhRi7rkVPCDd2SLjgsC6jDz6o7/IoUdX3O6U5wIoBa1NERNR7xhKfSMdJvQvjGgHzHiGkyrPo/4dljXh2TYw/SoM+SOe7/5bJNr6C2G/vCCIyvFiLYlx72i7Vz2Hd02RHPZ47kjhDzM5DnBt/ID3hwB4yO7QsLBzHNZG0zHvwRRJndE07bIw=="
	data, err := base64.StdEncoding.DecodeString(decText)
	assert.Nil(t, err)

	Decrypt, err := Decrypt(data, priKeyBytes)
	assert.Nil(t, err)

	if string(Decrypt) == text {
		fmt.Println("rsa parse success!")
	} else {
		fmt.Println("rsa parse Fail!")
	}

}
