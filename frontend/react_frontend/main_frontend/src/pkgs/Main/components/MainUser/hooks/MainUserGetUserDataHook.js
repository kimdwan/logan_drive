import { MainUserGetUserDataFunc } from "../functions"
import { useEffect, useState } from "react"

export const useMainUserGetUserDataHook = (computerNumber, setComputerNumber, go_backend_url, navigate) => {
    // 데이터 가져오기
    const [ detailData, setDetailData ] = useState("")
    const [ userProfile, setUserProfile ] = useState("")
  
    useEffect(() => {
      const detail_url = `${go_backend_url}/auth/get/detail`
      const profile_url = `${go_backend_url}/auth/get/profileimg`
      const getUserDataClass = new MainUserGetUserDataFunc(computerNumber, setComputerNumber, navigate)
      const getDataFunc = async (detail_url, profile_url, getUserDataClass) => {
  
        // 유저의 detail을 가져오는 로직
        const userDetail = await getUserDataClass.GetData(detail_url)
        
        // 유저의 프로필을 가져오는 로직
        const userProfile = await getUserDataClass.GetData(profile_url)
  
        if (userDetail && userProfile) {
          setDetailData(userDetail)
          
          // 이미지 만들어주기
          const img_url = `data:${userProfile["ImgType"]};charset=utf-8;base64,${userProfile["ImgBase64"]}`
          setUserProfile(img_url)
        }
  
      }
  
      if (computerNumber) {
        getDataFunc(detail_url, profile_url, getUserDataClass)
      }
  
    }, [ computerNumber, setComputerNumber, go_backend_url, navigate, setDetailData, setUserProfile ])

    return { detailData, userProfile }
}