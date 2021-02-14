# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source = [
  "./dist/kpkg_darwin_amd64/kpkg"]
bundle_id = "net.thespblog.kpkg"

apple_id {
  username = "spachava753@gmail.com"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Shashank Pachava"
}