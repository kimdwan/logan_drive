import { ContentsNavbarUserProfileFunction } from "../functions"

import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"

export const useContentsNavbarUserGetDataHook = (computerNumber, setComputerNumber ) => {
  const [ profileImg, setProfileImg ] = useState(undefined)
  const [ userDetail, setUserDetail ] = useState("")
  const navigate = useNavigate()

  useEffect(() => {
    const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL

    if ( computerNumber ) {

      const userDataClasse = new ContentsNavbarUserProfileFunction(computerNumber, setComputerNumber, navigate)
      
      // 유저의 프로필 이미지를 가져오는 로직
      const getUserProfileImgDataFunc = async (img_url) => {
      
        const img_Data = await userDataClasse.GetUserData(img_url)

        if (img_Data && img_Data["ImgBase64"] && img_Data["ImgType"]) {
          setProfileImg(`data:${img_Data["ImgType"]};charset=utf-8;base64,${img_Data["ImgBase64"]}`)
        }
      }

      // 유저의 닉네임과 이메일을 가져오는 로직
      const getUserDetailDataFunc = async (detail_url) => {

        const user_detail_data = await userDataClasse.GetUserData(detail_url)

        if (user_detail_data && user_detail_data["email"] && user_detail_data["nickname"]) {
          setUserDetail(`${user_detail_data["email"]}(${user_detail_data["nickname"]})`)
        }

      }


      // 데이터 가져오기
      getUserProfileImgDataFunc(`${go_backend_url}/auth/get/profileimg`)
      getUserDetailDataFunc(`${go_backend_url}/auth/get/detail`)

    }

  }, [ computerNumber, setComputerNumber, navigate ])


  return { profileImg, userDetail }
}