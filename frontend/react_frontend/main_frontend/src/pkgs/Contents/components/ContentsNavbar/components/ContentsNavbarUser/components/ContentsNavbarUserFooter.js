import { useCallback } from "react"
import { useNavigate } from "react-router-dom"
import { ContentsNavbarUserFooterFunc } from "../functions"

export const ContentNavbarUserFooter = ({ computerNumber, setComputerNumber }) => {
  const navigate = useNavigate()

  const clickLogoutBtn = useCallback((event) => {

    if (computerNumber && event.target.className === "contentNavbarUserFooterLogoutBtn") {

      const contentNavbarFooterClass = new ContentsNavbarUserFooterFunc(computerNumber, setComputerNumber, navigate)
      const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL 
      if (go_backend_url) {
        contentNavbarFooterClass.ClickLogoutBtn(`${go_backend_url}/auth/user/logout`)
      }

    }

  }, [ computerNumber, setComputerNumber, navigate ])


  return (
    <div className = "contentNavbarUserFooterContainer">

      {/* 실질적으로 값이 들어가는 장소 */}
      <div className = "contentNavbarUserFooterDivBox">
        
        {/* 로그아웃 버튼 */}
        <div className = "contentNavbarUserFooterLogoutDivBox">
          <button 
          className = "contentNavbarUserFooterLogoutBtn"
          onClick = { clickLogoutBtn }
          >
            로그아웃
          </button>
        </div>

      </div>

    </div>
  )
}