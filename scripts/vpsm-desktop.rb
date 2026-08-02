cask "vpsm-desktop" do
  version "1.1.7"
  sha256 :no_check

  url "https://github.com/devlopersabbir/vpcm/releases/download/v#{version}/vpsm-desktop-mac.dmg"
  name "VPSM Desktop"
  desc "VPS Manager - Remote Server Inventory & SSH Terminal Panel"
  homepage "https://github.com/devlopersabbir/vpcm"

  app "VPSM Desktop.app"

  zap trash: [
    "~/Library/Application Support/VPSM Desktop",
    "~/Library/Preferences/com.devlopersabbir.vpsm-desktop.plist",
    "~/Library/Saved Application State/com.devlopersabbir.vpsm-desktop.savedState",
  ]
end
