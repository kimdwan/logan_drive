import "./assets/css/ContentNavbarFriend.css"

import { ContentNavbarFriendList, ContentNavbarFriendTop } from "./components"

export const ContentNavbarFriend = ({ computerNumber, setComputerNumber }) => {

  return (
    <div className = "contentNavbarFriendContainer">
      
      {/* 타이틀 이름 */}
      <ContentNavbarFriendTop />

      {/* 친구 리스트가 보이는 장소 */}
      <ContentNavbarFriendList computerNumber = { computerNumber } setComputerNumber = { setComputerNumber } />

    </div>
  )
}