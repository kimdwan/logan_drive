import profileImg from "../assets/img/profileImg.png"

import { useContentNavbarFriendListHook, useContentNavbarFriendGetStatusHook } from "../hooks"

import { useNavigate } from "react-router-dom"

export const ContentNavbarFriendList = ({ computerNumber, setComputerNumber }) => {
  // navigate
  const navigate = useNavigate()

  // 친구의 리스트를 가져오는 함수
  const { friendList } = useContentNavbarFriendListHook( computerNumber, setComputerNumber, navigate )

  // 친구의 접속상태를 확인하는 상태 (객체 형태로 관리)
  const { friendStatus } = useContentNavbarFriendGetStatusHook( computerNumber )

  return (
    <div className = "contentNavbarFriendListContainer">

      {/* 친구의 리스트가 온다. */}
      <div className = "contentNavbarFriendListScheduleBox">

        {/* 추후에 즐겨찾기와 명패를 업로드 해야함 */}
        {
          friendList && friendList.length > 0 && 
          friendList.map((friendData, idx) => {
            // 친구 상태창 확인
            const statusData = friendStatus && friendStatus.length > 0 && friendStatus.find(status => status["friend_id"] === friendData["friend_id"])
            const statuses = statusData && statusData["status"]

            return (
              <div className = "contentNavbarFriendListScheduleData" key = {idx} data-friendid = { friendData["friend_id"] }>

                {/* 친구의 배경화면 */}
                <div className = "contentNavbarFriendListScheduleProfileDataBox">
                  <img 
                    className = "contentNavbarFriendListScheduleProfileImgData"
                    src = { friendData["friend_imgbase64"] && friendData["friend_imgtype"] ? `data:${friendData["friend_imgtype"]};charset=utf-8;base64,${friendData["friend_imgbase64"]}` : profileImg }
                    alt = "친구의 프로필"
                  />
                </div>

                {/* 친구의 이메일과 닉네임 */}
                <div className = "contentNavbarFriendListScheduleProfileEmailAndNickNameDataBox">
                  <div className = "contentNavbarFriendListScheduleProfileEmailAndNickNameValue">
                    {
                      friendData["friend_nickname"] && friendData["friend_email"] ? `${friendData["friend_nickname"]}(${friendData["friend_email"]})` : "찾을 수 없음"
                    }
                  </div>
                </div>

                {/* 친구의 접속상태 확인 */}
                <div className = "contentNavbarFriendListScheduleProfileConnectCheckBox">
                    {
                      statuses === 0 ? "⚪" : statuses === 1 ? "🟢" : statuses === 2 ? "🟡" : statuses === 3 ? "🔴" : statuses === 4 ? "❓" : "오류"
                    }
                </div>

              </div>
            )
          })
        }

      </div>

    </div>
  )
}