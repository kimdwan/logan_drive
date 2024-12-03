import "./assets/css/ContentNavbar.css"

import { ContentNavbarFriend, ContentNavbarTop, ContentNavbarUser } from "./components"
import { MainContext } from "../../../../App"

import { useContext } from "react"

export const ContentNavbar = () => {
  const { computerNumber, setComputerNumber } = useContext(MainContext)

  return (
    <div className = "contentsNavbarContainer">

      {/* 유저의 정보가 보이는 컴퍼넌트 */}
      <ContentNavbarUser computerNumber = {computerNumber} setComputerNumber = {setComputerNumber} />

      {/* 가장 상단에 위치한 컴퍼넌트 */}
      <ContentNavbarTop />

      {/* 친구 창이 보이는 컴퍼넌트 */}
      <ContentNavbarFriend />

    </div>
  )
}