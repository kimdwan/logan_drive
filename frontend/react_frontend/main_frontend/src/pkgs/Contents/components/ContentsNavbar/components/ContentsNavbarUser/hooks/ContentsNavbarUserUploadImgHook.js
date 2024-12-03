import { ContentsNavbarUploadImgFunction } from "../functions"

import { useCallback, useState } from "react"
import { useForm } from "react-hook-form"
import { useNavigate } from "react-router-dom"

export const useContentsNavbarUploadImgHook = (computerNumber, setComputerNumber) => {
  const { register, setValue, handleSubmit } = useForm()
  
  const [ userImgData, setUserImgData ] = useState(undefined)

  const navigate = useNavigate()
  const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL

  const onSubmit = async (data) => {
    
    try {
      // formData로 변경하기
      const formData = new FormData()

      const url_profile_data = data["user_profile_img"]

      if ( url_profile_data["name"] ) {
        formData.append("user_profile_img", url_profile_data)

        if (computerNumber) {
          const uploadClass = new ContentsNavbarUploadImgFunction(computerNumber, setComputerNumber, navigate)

          const message = await uploadClass.UploadUserProfile(`${go_backend_url}/auth/user/upload/profile`, formData)

          if (message) {
            alert(message)
            window.location.reload()
          }
        }
      } else {
        alert("파일을 업로드 해주셔야 합니다")
        throw new Error("파일을 업로드 안함")
      }

    } catch (err) {
      throw err
    }

  }

  const clickUploadImgBtn = useCallback((event) => {

    if (event.target.id === "contentNavbarUserUploadImgConveyInterectBox") {
      try { 
        const data = event.target.files[0]

        // 데이터 크기 검정 
        const data_size = data["size"]
        if (data_size > (10 * 1024 * 1024)) {
          alert("데이터 크기는 10mb를 넘을수 없습니다")
          throw new Error("데이터 크기가 10mb를 넘음")
        }

        // 데이터 타입 검정
        const data_name_list = String(data["name"]).split(".")
        let istypeAllowed = false
        const data_type = data_name_list[data_name_list.length - 1].toLowerCase()
        const allow_data_types = Array.from(["webp", "png", "jpg", "jpeg"])
        for (let i = 0; i < allow_data_types.length; i++) {
          const allow_data_type = allow_data_types[i]
          if (data_type === allow_data_type) {
            istypeAllowed = true  
            break
          }
        }

        if (!istypeAllowed) {
          alert("올릴수 있는 이미지 타입은 png, jpg, jpeg, webp입니다")
          throw new Error("허용하지 않은 이미지 타입")
        }

        setValue("user_profile_img", data)
        setUserImgData(data["name"])

      } catch(err) {
        throw err
      }
    }

  },[ setValue, setUserImgData ])

  return { register, handleSubmit, onSubmit, clickUploadImgBtn, userImgData }
}