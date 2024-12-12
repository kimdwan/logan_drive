import { useContentNavbarTopMidGetTitleNameHook } from "../hooks"

export const ContentNavbarTopMid = () => {
  const { titleName } = useContentNavbarTopMidGetTitleNameHook()

  return (
    <div className = "contentNavbarTopMidContainer">

      {/* 타이틀이 들어가는 장소 */}
      <div className = "contentNavbarTopMidTitleDiv">
        <h2 className = "contentNavbarTopMidTitleValue">
          {
            titleName
          }
        </h2>
      </div>

    </div>
  )
}