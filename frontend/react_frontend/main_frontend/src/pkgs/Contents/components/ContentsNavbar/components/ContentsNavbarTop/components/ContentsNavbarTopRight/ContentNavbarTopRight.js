import "./assets/css/ContentNavbarTopRight.css"

import { ContentNavbarTopRightChannelListEdition } from "./components"
import { useContentNavbarTopRightGetUrlNameHook } from "./hooks"

export const ContentNavbarTopRight = () => {
  const { urlName } = useContentNavbarTopRightGetUrlNameHook()

  return (
    <div className = "contentNavbarTopRightContainer">
      {
        urlName === "channellist" ? <ContentNavbarTopRightChannelListEdition /> : <></>
      }
    </div>
  )
}