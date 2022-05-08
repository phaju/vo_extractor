package main

import (
	"encoding/json"
	"log"
	"os"
  "fmt"
	"path/filepath"
	"strings"
  "sort"
)

func main() {
    currentDirectory, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    vo_count, vo_list := iterate(currentDirectory)

    f, err := os.Create("Extracted_VO.txt")
    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    f.WriteString("Total VOs found: " + fmt.Sprint(vo_count) + "\n")
    f.Write([]byte(vo_list))
}

func iterate(path string) (int, string) {
    var vo_list []string
    filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Fatalf(err.Error())
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
                json_data := dat["json"].(map[string]interface{})
                steps_array := json_data["steps"].([]interface{})
                for step_number := range steps_array {
                  current_step := steps_array[step_number].(map[string]interface{})
                  info_vo := get_vo(current_step)
                  if string(info_vo) != "" {
                    vo_list = append(vo_list, string(info_vo)) 
                  } else {
                    question_object := current_step["question"].(map[string]interface{})
                    question_vo := get_vo(question_object)
                    if string(question_vo) != "" {
                      vo_list = append(vo_list, string(question_vo))
                    }
                    feedback_object := current_step["feedBack"].([]interface{})
                    for feedback_index := range feedback_object {
                      current_feedback := feedback_object[feedback_index].(map[string]interface{})
                      feedback_vo := get_vo(current_feedback)
                      if string(feedback_vo) != "" {
                        vo_list = append(vo_list, string(feedback_vo)) 
                      }
                    }
                  }
                }
              } 
            }
          }
        }
        return nil
    })
    sort.Strings(vo_list)
    return len(vo_list), strings.Join(vo_list, "\n")
}

func get_vo(vo_obj map[string]interface{}) []byte {
    var vo string
    if vo_obj["text"] != nil {
      if audio_location := vo_obj["audioFile"]; audio_location != nil {
        if path_content := strings.Split(audio_location.(string), "/"); len(path_content) == 3 {
          vo = path_content[2]
          vo = vo + " : "
        }
      }
      vo = vo + vo_obj["text"].(string)
    }
    return []byte(vo)
}
