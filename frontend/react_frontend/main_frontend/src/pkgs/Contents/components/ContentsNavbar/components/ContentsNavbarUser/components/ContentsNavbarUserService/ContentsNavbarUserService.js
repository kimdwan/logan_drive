import { ContentNavbarUserServiceChannelList } from "./components"
import { useContentNavbarUserServiceReedemLocationHook } from "./hooks"

export const ContentNavbarUserService = () => {
  const { urlTitle } = useContentNavbarUserServiceReedemLocationHook()
  
  return (
    <div className = "contentNavbarUserServiceContainer">
      {
        urlTitle === "channellist" ? <ContentNavbarUserServiceChannelList /> : <></>
      }
    </div>
  )
}