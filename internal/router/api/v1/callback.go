package v1

import (
	"firstProject/internal/app"
	domainDrive "firstProject/internal/domain/drive"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
)

func Callback(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		events, err := app.LineBotClient.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Println(err)
				c.JSON(http.StatusBadRequest, err)
			} else {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, err)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if message.Text == "login" {
						lineID := event.Source.UserID
						authURL := app.DriveService.LoginURL(ctx, lineID)
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(authURL)).Do(); err != nil {
							log.Println(err)
						}
						return
					}
					if message.Text == "list" {
						lineID := event.Source.UserID
						res, err := app.DriveService.ListFiles(ctx, lineID)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintln(res))).Do(); err != nil {
							log.Println(err)
						}
						return
					}
					if message.Text == "list folder" {
						lineID := event.Source.UserID
						res, err := app.DriveService.ListMyDriveFolders(ctx, lineID)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintln(res))).Do(); err != nil {
							log.Println(err)
						}
						return
					}
					if message.Text == "list shared" {
						lineID := event.Source.UserID
						res, err := app.DriveService.ListSharedFolders(ctx, lineID)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintln(res))).Do(); err != nil {
							log.Println(err)
						}
						return
					}
					if message.Text == "test" {
						lineID := event.Source.UserID
						res, err := app.DriveService.TestFolderCarousel(ctx, lineID)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err := app.LineBotClient.ReplyMessage(
							event.ReplyToken,
							linebot.NewFlexMessage("測試Flex Carousel", res.CarouselContainer),
						).Do(); err != nil {
							log.Println(err)
							return
						}
					}
					if message.Text == "mydrive" {
						lineID := event.Source.UserID
						res, err := app.DriveService.ListFolderCarousel(ctx, lineID, domainDrive.PersonalFolder)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err := app.LineBotClient.ReplyMessage(
							event.ReplyToken,
							linebot.NewFlexMessage("測試Flex Carousel", res.CarouselContainer),
						).Do(); err != nil {
							log.Println(err)
							return
						}
					}
					if message.Text == "shared" {
						lineID := event.Source.UserID
						res, err := app.DriveService.ListFolderCarousel(ctx, lineID, domainDrive.SharedFolder)
						if err != nil {
							log.Println(err)
							return
						}
						if _, err := app.LineBotClient.ReplyMessage(
							event.ReplyToken,
							linebot.NewFlexMessage("測試Flex Carousel", res.CarouselContainer),
						).Do(); err != nil {
							log.Println(err)
							return
						}
					}
					if message.Text == "flex carousel" {
						contents := &linebot.CarouselContainer{
							Type: linebot.FlexContainerTypeCarousel,
							Contents: []*linebot.BubbleContainer{
								{
									Type: linebot.FlexContainerTypeBubble,
									Body: &linebot.BoxComponent{
										Type:   linebot.FlexComponentTypeBox,
										Layout: linebot.FlexBoxLayoutTypeVertical,
										Contents: []linebot.FlexComponent{
											&linebot.TextComponent{
												Type:   linebot.FlexComponentTypeText,
												Text:   "FOLDER",
												Weight: linebot.FlexTextWeightTypeBold,
												Color:  "#1DB446",
												Size:   linebot.FlexTextSizeTypeSm,
											},
											&linebot.TextComponent{
												Type:   linebot.FlexComponentTypeText,
												Text:   "Folder Name",
												Weight: linebot.FlexTextWeightTypeBold,
												Size:   linebot.FlexTextSizeTypeXxl,
												Margin: linebot.FlexComponentMarginTypeMd,
											},
											&linebot.TextComponent{
												Type:  linebot.FlexComponentTypeText,
												Text:  "/path/to/floder",
												Size:  linebot.FlexTextSizeTypeXs,
												Color: "#aaaaaa",
												Wrap:  true,
											},
											&linebot.SeparatorComponent{
												Type:   linebot.FlexComponentTypeSeparator,
												Margin: linebot.FlexComponentMarginTypeXxl,
											},
											&linebot.BoxComponent{
												Type:    linebot.FlexComponentTypeBox,
												Layout:  linebot.FlexBoxLayoutTypeVertical,
												Margin:  linebot.FlexComponentMarginTypeXxl,
												Spacing: linebot.FlexComponentSpacingTypeSm,
												Contents: []linebot.FlexComponent{
													&linebot.BoxComponent{
														Type:   linebot.FlexComponentTypeBox,
														Layout: linebot.FlexBoxLayoutTypeHorizontal,
														Contents: []linebot.FlexComponent{
															&linebot.TextComponent{
																Type:       linebot.FlexComponentTypeText,
																Text:       "Folder1",
																Size:       linebot.FlexTextSizeTypeSm,
																Color:      "#555555",
																Decoration: linebot.FlexTextDecorationTypeUnderline,
																MaxLines:   linebot.IntPtr(25),
																Align:      linebot.FlexComponentAlignTypeStart,
																Margin:     linebot.FlexComponentMarginTypeNone,
																Gravity:    linebot.FlexComponentGravityTypeCenter,
																Flex:       linebot.IntPtr(0),
															},
															&linebot.FillerComponent{
																Type: linebot.FlexComponentTypeFiller,
															},
															&linebot.ButtonComponent{
																Type: linebot.FlexComponentTypeButton,
																Action: &linebot.PostbackAction{
																	Label:       "進入資料夾",
																	Data:        "folderid1",
																	DisplayText: "進入Folder1",
																},
																Style:      linebot.FlexButtonStyleTypeLink,
																Height:     linebot.FlexButtonHeightTypeSm,
																Gravity:    linebot.FlexComponentGravityTypeCenter,
																Flex:       linebot.IntPtr(0),
																AdjustMode: linebot.FlexComponentAdjustModeTypeShrinkToFit,
															},
														},
													},
													&linebot.BoxComponent{
														Type:   linebot.FlexComponentTypeBox,
														Layout: linebot.FlexBoxLayoutTypeHorizontal,
														Contents: []linebot.FlexComponent{
															&linebot.TextComponent{
																Type:       linebot.FlexComponentTypeText,
																Text:       "Folder2",
																Size:       linebot.FlexTextSizeTypeSm,
																Color:      "#555555",
																Decoration: linebot.FlexTextDecorationTypeUnderline,
																MaxLines:   linebot.IntPtr(25),
																Align:      linebot.FlexComponentAlignTypeStart,
																Margin:     linebot.FlexComponentMarginTypeNone,
																Gravity:    linebot.FlexComponentGravityTypeCenter,
																Flex:       linebot.IntPtr(0),
															},
															&linebot.FillerComponent{
																Type: linebot.FlexComponentTypeFiller,
															},
															&linebot.ButtonComponent{
																Type: linebot.FlexComponentTypeButton,
																Action: &linebot.PostbackAction{
																	Label:       "進入資料夾",
																	Data:        "folderid2",
																	DisplayText: "進入Folder2",
																},
																Style:      linebot.FlexButtonStyleTypeLink,
																Height:     linebot.FlexButtonHeightTypeSm,
																Gravity:    linebot.FlexComponentGravityTypeCenter,
																Flex:       linebot.IntPtr(0),
																AdjustMode: linebot.FlexComponentAdjustModeTypeShrinkToFit,
															},
														},
													},
												},
											},
											// Separator
											&linebot.SeparatorComponent{
												Margin: linebot.FlexComponentMarginTypeXxl,
											},
											// Files
											&linebot.BoxComponent{
												Type:   linebot.FlexComponentTypeBox,
												Layout: linebot.FlexBoxLayoutTypeHorizontal,
												Margin: linebot.FlexComponentMarginTypeXxl,
												Contents: []linebot.FlexComponent{
													&linebot.TextComponent{
														Type:  linebot.FlexComponentTypeText,
														Text:  "Total Files",
														Size:  linebot.FlexTextSizeTypeSm,
														Color: "#555555",
														Flex:  linebot.IntPtr(0),
													},
													&linebot.TextComponent{
														Type:  linebot.FlexComponentTypeText,
														Text:  "3",
														Size:  linebot.FlexTextSizeTypeSm,
														Color: "#111111",
														Align: linebot.FlexComponentAlignTypeEnd,
													},
												},
											},
											&linebot.BoxComponent{
												Type:   linebot.FlexComponentTypeBox,
												Layout: linebot.FlexBoxLayoutTypeHorizontal,
												Contents: []linebot.FlexComponent{
													&linebot.TextComponent{
														Type:  linebot.FlexComponentTypeText,
														Text:  "FILE1",
														Size:  linebot.FlexTextSizeTypeSm,
														Color: "#555555",
													},
												},
											},
											&linebot.BoxComponent{
												Type:   linebot.FlexComponentTypeBox,
												Layout: linebot.FlexBoxLayoutTypeHorizontal,
												Contents: []linebot.FlexComponent{
													&linebot.TextComponent{
														Type:  linebot.FlexComponentTypeText,
														Text:  "FILE2",
														Size:  linebot.FlexTextSizeTypeSm,
														Color: "#555555",
													},
												},
											},
											&linebot.BoxComponent{
												Type:   linebot.FlexComponentTypeBox,
												Layout: linebot.FlexBoxLayoutTypeHorizontal,
												Contents: []linebot.FlexComponent{
													&linebot.TextComponent{
														Type:  linebot.FlexComponentTypeText,
														Text:  "FILE3",
														Size:  linebot.FlexTextSizeTypeSm,
														Color: "#555555",
													},
												},
											},
											// Separator
											&linebot.SeparatorComponent{
												Margin: linebot.FlexComponentMarginTypeXxl,
											},
											// Button
											&linebot.BoxComponent{
												Type:   linebot.FlexComponentTypeBox,
												Layout: linebot.FlexBoxLayoutTypeHorizontal,
												Margin: linebot.FlexComponentMarginTypeMd,
												Contents: []linebot.FlexComponent{
													&linebot.ButtonComponent{
														Type: linebot.FlexComponentTypeButton,
														Action: &linebot.PostbackAction{
															Label:       "設為上傳資料夾",
															Data:        "folderid",
															DisplayText: "設為上傳資料夾",
														},
														Style:      linebot.FlexButtonStyleTypePrimary,
														AdjustMode: linebot.FlexComponentAdjustModeTypeShrinkToFit,
													},
												},
											},
										},
									},
									Styles: &linebot.BubbleStyle{
										Footer: &linebot.BlockStyle{
											Separator: true,
										},
									},
								},
							},
						}
						if _, err := app.LineBotClient.ReplyMessage(
							event.ReplyToken,
							linebot.NewFlexMessage("Flex message alt text", contents),
						).Do(); err != nil {
							log.Println(err)
							return
						}
					}
					samplePK, err := app.SampleService.Sample(ctx, message.Text)
					if err != nil {
						log.Println(err)
						return
					}
					if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(samplePK)).Do(); err != nil {
						log.Println(err)
					}

				case *linebot.FileMessage:
					lineID := event.Source.UserID
					content, err := app.LineBotClient.GetMessageContent(message.ID).Do()
					if err != nil {
						log.Println(err)
						return
					}

					app.DriveService.UploadFile(ctx, lineID, message.FileName, content.Content)
					if err != nil {
						log.Println(err)
						return
					}
					if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("成功上傳: "+message.FileName)).Do(); err != nil {
						log.Println(err)
					}

				case *linebot.ImageMessage:
					lineID := event.Source.UserID
					// 使用 Line Bot API 獲取圖片內容
					content, err := app.LineBotClient.GetMessageContent(message.ID).Do()
					if err != nil {
						log.Printf("Failed to get image content: %s", err)
						return
					}
					app.DriveService.UploadFile(ctx, lineID, message.ID, content.Content)
					if err != nil {
						log.Println(err)
						return
					}
					if _, err = app.LineBotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("成功上傳: "+message.ID)).Do(); err != nil {
						log.Println(err)
					}
				}
			}
		}

	}

}
