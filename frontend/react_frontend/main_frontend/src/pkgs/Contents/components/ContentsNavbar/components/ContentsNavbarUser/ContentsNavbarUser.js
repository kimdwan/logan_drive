import "./assets/css/ContentNavbarUser.css"

import { ContentNavbarUserFooter, ContentNavbarUserLogo, ContentNavbarUserProfile, ContentNavbarUserService, ContentNavbarUserUploadImg } from "./components"

export const ContentNavbarUser = ({ computerNumber, setComputerNumber }) => {
  return (
    <div className = "contentNavbarUserContainer">
      
      {/* 로고가 위치함 */}
      <ContentNavbarUserLogo />

      {/* 프로필 정보가 위치함 */}
      <ContentNavbarUserProfile computerNumber = {computerNumber} setComputerNumber = {setComputerNumber} />

      {/* 유저가 이미지를 업로드 할 수 있음 */}
      <ContentNavbarUserUploadImg />

      {/* 유저가 서비스 메뉴를 선택할 수 있음 */}
      <ContentNavbarUserService />

      {/* 아래쪽을 담당하고 있음 */}
      <ContentNavbarUserFooter />

    </div>
  )
}