import subMainLogo from "../assets/img/subMainLogo.webp"

export const MainLoginLogo = () => {
  return (
    <div className = "mainLoginLogo">
      <img className = "mainLoginLogoImage" src = { subMainLogo } alt = "보조 로고" />
    </div>
  )
}