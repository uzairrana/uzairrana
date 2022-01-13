package Authentication

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/speps/go-hashids"
)

type EmailVerification struct {
	EMAIL string `json:"EMAIL"`
}

type code struct {
	Code string
}

type EmailVerificationCode struct {
	EmailCode string `json:"EmailCode"`
}

type IsActive struct {
	IsActive string
}

func VerifyEmail(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	smtpHost := "smtp.gmail.com"

	userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	from := "noor.fatima5675@gmail.com" //"argontechgaming@gmail.com"
	logger.Info("====================================VerifyEmail Start=====================================: " + payload)
	smtpPort := "587"
	emailVerification := &EmailVerification{}
	err := json.Unmarshal([]byte(payload), &emailVerification)
	if err != nil {
		//logger.Info(err)
		//logger.Info("error unmarshal payload to emailVerification -- VerifyViaEmail -- ", errors.New("error please try again 301"))
		return "", errors.New("Parsing error please try again 301")
	} else {
		tempEmail := emailVerification.EMAIL
		logger.Info("===================================Email is: :: %v " + tempEmail) // yeh kia hae  ?
		//Email validation start
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		//Email validation end
		if re.MatchString(tempEmail) != true {
			//fmt.Info(re)
			return "", errors.New("error please try again 301")
		} else {
			//Random code generation start
			hd := hashids.NewData()
			hd.Salt = "this is my salt"
			h, _ := hashids.NewWithData(hd)
			timestamp := time.Now().UnixNano()
			t := int(timestamp)
			idd, _ := h.Encode([]int{t})
			idd = idd[0:6]
			var id = idd //"12345"
			logger.Info(id)
			message := "Please enter this code for verification in the Poker Game.  ' " + id + " '\n Thanks"

			//by M.Talha Saleem

			auth := smtp.PlainAuth("", from, "Go@Fast!2", smtpHost)

			//var to []string
			to := emailVerification.EMAIL //append(to, emailVerification.EMAIL)
			subject := "Signup Code"
			body := message + "  " + time.Now().String()
			msg := composeMimeMail(to, from, subject, body)
			//err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
			err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)

			if err != nil {
				fmt.Println(err)

				logger.Info("Verification code has been not sent")

				return "Verification code has been not sent", nil
			}
			fmt.Println("Email Sent Successfully!")

			logger.Info("Email Sent Successfully!")

			// by M.Talha Saleem

			//Random code generation end

			// //Email send start
			// //auth := smtp.PlainAuth("", "mobeenarchi@gmail.com", "mobeenarchi2", "smtp.gmail.com")
			// auth := smtp.PlainAuth("", "verify@homegamepoker.io", "3mCL*EJh", "smtp.gmail.com")

			// // Connect to the server, authenticate, set the sender and recipient,
			// // and send the email all in one step.
			// to := []string{"mobeenarchi@gmail.com"}
			// //logger.Info("===========   to")
			// msg := []byte("To: " + strings.Join(to, "") + "  \r\n" +
			// 	"Subject: Home Game Poker - Verification Code\r\n" +
			// 	"\r\n" +
			// 	"Confirmation Code: '" + id + "' \r\n" +
			// 	"Please verify your Home Game Poker account by entering this code on the Confirmation screen inside the app.\r\n")
			// logger.Info("===========   msg")
			// //err := smtp.SendMail("smtp.gmail.com:587", auth, "mobeenarchi@gmail.com", to, msg)
			// err := smtp.SendMail("smtp.gmail.com:587", auth, "verify@homegamepoker.io", to, msg)
			// if err != nil {
			// 	log.Fatal(err)
			// }

			//Email send start
			logger.Info("===================================UUID is: :: %v " + id)
			email_code := code{id} // naming baad mae theek kr rlena error ko read krna alzmi ana chaye
			uuid, _ := json.Marshal(email_code)
			//Email code save in the Storage
			objectIds := []*runtime.StorageWrite{
				&runtime.StorageWrite{
					Collection:      "EmailVerification",
					Key:             "VerificationCode",
					UserID:          userID,
					Value:           string(uuid),
					PermissionRead:  1,
					PermissionWrite: 0,
				},
			}

			if _, err := nk.StorageWrite(ctx, objectIds); err != nil {
				//logger.Info("Error writing verification code in DB -- EmailVerification -- ", err)
			}

			//User Active/InActive status write start
			isActive := IsActive{"False"}
			status, _ := json.Marshal(isActive)
			//User status save in the Storage
			objectIds2 := []*runtime.StorageWrite{
				&runtime.StorageWrite{
					Collection:      "UserActiveState",
					Key:             "isActive",
					UserID:          userID,
					Value:           string(status),
					PermissionRead:  1,
					PermissionWrite: 0,
				},
			}

			if _, err := nk.StorageWrite(ctx, objectIds2); err != nil {
				//logger.Info("Error writing user isActive state in DB -- UserActiveState -- ", err)
			}
			//User Active/InActive State write end
		}
	}

	logger.Info("====================================VerifyEmail End=====================================")

	//  idr humhyay email ae ga hum bhi uski regex ko verify kry gae agr galat hogi toph error return kr dae gae agae kuch nhi hoga
	//agr theek email regex hogi toh aik code hgenerate krwa kr us email mae send kry gae
	// or wo code player kai storag3e mae rakhwa dae gae

	// or player ko notify kry gae kae wo account code sae verify kry

	return "EmailSent", nil
}

func VerifyCode(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	logger.Info("====================================VerifyCode Start=====================================: " + payload)
	emailVerificationCode := &EmailVerificationCode{}
	err := json.Unmarshal([]byte(payload), &emailVerificationCode)
	if err != nil {
		//	logger.Info(err)
		//logger.Info("error unmarshal payload to verifiEmaiCode -- VerifyCode -- ", errors.New("error please try again 301"))
		return "", errors.New("VerifyCode Parsing error please try again 301")
	} else {
		objectIds := []*runtime.StorageRead{
			&runtime.StorageRead{
				Collection: "EmailVerification",
				Key:        "VerificationCode",
				UserID:     userID,
			},
		}
		// all records fetched
		records, err := nk.StorageRead(ctx, objectIds)

		if err != nil {
			//logger.Info("error reading from UserActiveState collection -- VerifyCode -- ", err)
			return "", errors.New("please try again error while getting your " + emailVerificationCode.EmailCode + " Collection")
		}

		if len(records) >= 1 {
			// status := false
			for _, record := range records {
				code := &code{}
				err := json.Unmarshal([]byte(record.Value), &code)
				if err != nil {
					//logger.Info(err)
					//logger.Info("error unmarshal EmailVerification Collection -- ", errors.New("error please try again 301"))
				} else {
					logger.Info("Record value is: ", code.Code)
					logger.Info("Record value unmarhsal is: ", emailVerificationCode.EmailCode)
					if code.Code == emailVerificationCode.EmailCode {
						// ----- User account verified via email ---- //

						//User Active/InActive State write start
						isActive := IsActive{"True"}
						state, _ := json.Marshal(isActive)
						//Email code save in the Storage
						objectIds2 := []*runtime.StorageWrite{
							&runtime.StorageWrite{
								Collection:      "UserActiveState",
								Key:             "isActive",
								UserID:          userID,
								Value:           string(state),
								PermissionRead:  1,
								PermissionWrite: 0,
							},
						}

						if _, err := nk.StorageWrite(ctx, objectIds2); err != nil {
							//logger.Info("Error writing user isActive state in DB -- UserActiveState -- ", err)
						}
						//User Active/InActive State write end
						logger.Info("----------------------------------------------------------")
						logger.Info(" Email verified, ", emailVerificationCode.EmailCode)
						logger.Info("----------------------------------------------------------")
						AssignMetaData(ctx, logger, db, nk, userID)
						AssignUserMatches(ctx, logger, nk, userID)
						return "Email verified", nil

					} else {
						logger.Info("Coode Doesn't match")
					}
				}
			}
		} else {
			// ----- User account not verified via email ---- //
			logger.Info("----------------------------------------------------------")
			logger.Info(" Email not verified, ", emailVerificationCode.EmailCode)
			logger.Info("----------------------------------------------------------")
			return "NotVerified", errors.New("Email not verified")
		}
	}
	logger.Info("====================================VerifyCode End=====================================: " + payload)
	//agly rpc mae client humhy rpc kae htorufgh  code dae ga
	// hum us code ko player account storage mae save code kae sath match kry gae
	// agar match ho jaye ga toh client ko notify kr dae gae wo account create kr dae using thro7ugh true  bool
	// or match krny kae baad uscode to storage sae delete krwa dae gae  to save a memory

	// agr match nhi howa toh client ko error send kry gae code soes not match

	return "", nil
}

func GetCurrentDateTime(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (string, error) {
	logger.Info("====================================GetCurrentDateTime Start=====================================: ")
	var currentTime int64 = time.Now().Unix()

	logger.Info("====================================GetCurrentDateTime End=====================================: ")
	return strconv.Itoa(int(currentTime)), nil
}

// M.Talha Saleem
func formatEmailAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}
	return e.String()
}

func encodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Address: str}
	return strings.Trim(addr.String(), " <>")
}

func composeMimeMail(to string, from string, subject string, body string) []byte {
	header := make(map[string]string)
	header["From"] = formatEmailAddress(from)
	header["To"] = formatEmailAddress(to)
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return []byte(message)
}
