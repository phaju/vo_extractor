package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

  err = os.Mkdir("Extracted_VO", 0755)
	if err != nil {
		log.Fatal(err)
	}

	iterate(currentDirectory)
}

func iterate(path string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err.Error())
		}
		if info.IsDir() == false {
			if filepath.Ext(info.Name()) == ".json" {
				file_content, err := os.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}

				var dat map[string]interface{}

				defer func() {
					if r := recover(); r != nil {
						//log.Println("Recovered. Error:\n", r)
					}
				}()

				if err := json.Unmarshal(file_content, &dat); err != nil {
					panic(err)
				}

				if dat["_name"] != nil {
					if dat["_name"] == "storyBoard" {
						var vo_list []string
						var mid string
						json_data := dat["json"].(map[string]interface{})
						mid = json_data["gameName"].(string)

						f, err := os.Create("Extracted_VO/" + mid + "_VO_list.csv")
						if err != nil {
							log.Fatal(err)
						}

						defer f.Close()

						game_text_array, ok := json_data["gameTexts"].(map[string]interface{})
						if ok {
							for key, value := range game_text_array {
								info_text := ""
								if value != "" && key != "" {
									info_text = fmt.Sprintf("%v,\"%v\"", key, value)
								}
								vo_list = append(vo_list, info_text)
							}
						}

						steps_array := json_data["steps"].([]interface{})
						for step_number := range steps_array {
							// Check for single info step
							current_step, ok := steps_array[step_number].(map[string]interface{})
							if ok {
								info_vo := get_vo(current_step)
								if string(info_vo) != "" {
									vo_list = append(vo_list, string(info_vo))
								}
							}
							// Check for multiple info
							info_object_array, ok := current_step["info"].([]interface{})
							if ok {
								for index := range info_object_array {
									current_feedback, ok := info_object_array[index].(map[string]interface{})
									if ok {
										feedback_vo := get_vo(current_feedback)
										if string(feedback_vo) != "" {
											vo_list = append(vo_list, string(feedback_vo))
										}
									}
								}
							}
							// Check for question step
							question_object, ok := current_step["question"].(map[string]interface{})
							if ok {
								question_vo := get_vo(question_object)
								if string(question_vo) != "" {
									vo_list = append(vo_list, string(question_vo))
								}
							}
							// Check for multiple question
							question_object_array, ok := current_step["question"].([]interface{})
							if ok {
								for index := range question_object_array {
									current_feedback, ok := question_object_array[index].(map[string]interface{})
									if ok {
										feedback_vo := get_vo(current_feedback)
										if string(feedback_vo) != "" {
											vo_list = append(vo_list, string(feedback_vo))
										}
									}
								}
							}
							// Check for question feedbacks
							feedback_object, ok := current_step["feedBack"].([]interface{})
							if ok {
								for index := range feedback_object {
									current_feedback, ok := feedback_object[index].(map[string]interface{})
									if ok {
										feedback_vo := get_vo(current_feedback)
										if string(feedback_vo) != "" {
											vo_list = append(vo_list, string(feedback_vo))
										}
									}
								}
							}
						}
						sort.Strings(vo_list)

						f.WriteString("Total VOs found," + fmt.Sprint(len(vo_list)) + "\n")
						f.WriteString("VO ID, English VO Text\n")
						f.Write([]byte(strings.Join(vo_list, "\n")))
					}
				}
			}
		}
		return nil
	})
}

func get_vo(vo_obj map[string]interface{}) []byte {
	deletion_list := []string{"<b>", "</b>", "<i>", "</i>"}
	var vo string
	if vo_obj["text"] != nil {
		if audio_location := vo_obj["audioFile"]; audio_location != nil {
			if path_content := strings.Split(audio_location.(string), "/"); len(path_content) > 2 {
				vo = path_content[len(path_content)-1]
			}
		}
		vo = vo + ",\"" + vo_obj["text"].(string) + "\""
		for index := range deletion_list {
			vo = strings.ReplaceAll(vo, deletion_list[index], "")
		}
	}
	return []byte(vo)
}
