import "./assets/css/ContentNavbarTop.css"
import { ContentNavbarTopLeft, ContentNavbarTopMid, ContentNavbarTopRight } from "./components"

export const ContentNavbarTop = () => {
  return (
    <div className = "contentNavbarTopContainer">
      
      {/* 네브바 탑에 왼쪽 담당 */}
      <ContentNavbarTopLeft />

      {/* 네브바 탑에 중간 담당 */}
      <ContentNavbarTopMid />

      {/* 네브바 탑에 오른쪽 담당 */}
      <ContentNavbarTopRight />

    </div>
  )
}