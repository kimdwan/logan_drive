import { useMainUserFooterBtnHook } from "../hooks"

import { useNavigate } from "react-router-dom"

export const MainUserFooter = ({ computerNumber, setComputerNumber }) => {

  // 기본 설정 들
  const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL
  const navigate = useNavigate()

  // 함수가 진행되는 로직
  const { clickFooterBtn } = useMainUserFooterBtnHook(computerNumber, setComputerNumber, go_backend_url, navigate)

  return (
    <div className = "mainUserFooterContainer">
      
      {/* 본격적으로 존재하는 아래 버튼 박스 */}
      <div className = "mainUserFooterDivBox">

        {/* 메인화면으로 가는 박스 */}
        <div className = "mainUserFooterGoMainDivBox">
          <button className = "mainUserFooterGoMainBtn" onClick = { clickFooterBtn } style={ { cursor : "pointer" } }>
            메인
          </button>
        </div>  

        {/* 로그아웃 버튼 */}
        <div className = "mainUserFooterLogoutBtnDivBox">
          <button className = "mainUserFooterLogoutBtn" onClick = { clickFooterBtn } style={ { cursor : "pointer" } }>
            LOGOUT
          </button>
        </div>

      </div>

    </div>
  )
}