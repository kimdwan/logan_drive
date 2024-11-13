import baseProfileImg from "../assets/img/baseProfileImg.png"

import { useMainUserGetUserDataHook } from "../hooks"

import { useNavigate } from "react-router-dom"

export const MainUserProfile = ({ computerNumber, setComputerNumber }) => {
  // 기본 설정
  const navigate = useNavigate()
  const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL

  // 유저 데이터 가져오기
  const { detailData, userProfile } = useMainUserGetUserDataHook(computerNumber, setComputerNumber, go_backend_url, navigate)

  return (
    <div className = "mainUserProfileContainer">
      
      {/* 프로필 이미지가 들어가는 컴퍼넌트 */}
      <div className = "mainUserProfileImgDivBox">
        <img className = "mainUserProfileImg" src = { userProfile || baseProfileImg } alt = "기본 프로필 이미지" style = {{ cursor : "pointer" }}/>
      </div>

      {/* 닉네임과 이메일이 존재하는 컴퍼넌트 */}
      <div className = "mainUserProfileEmailAndNicknameDivBox">
        
        {/* 닉네임 */}
        <div className = "mainUserProfileNickNameDivBox">
          <h2 className = "mainUserProfileNickNameValue">{ detailData["nickname"] }</h2>
        </div>

        {/* 이메일 */}
        <div className = "mainUserProfileEmailDivBox">
          <h4 className = "mainUserProfileEmailValue">{ detailData["email"] }</h4>
        </div>
      </div>

    </div>
  )
}