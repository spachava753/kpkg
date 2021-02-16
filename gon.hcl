source = ["./dist/kpkg-macos_darwin_amd64/kpkg"]
bundle_id = "net.thespblog.kpkg"

apple_id {
  username = "spachava753@gmail.com"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "8195B209BB79A5FF33C728DC1BBA3B5F30D5869D"
}

notarize {
  path = "./dist/kpkg-macos_darwin_amd64/kpkg"
  bundle_id = "net.thespblog.kpkg"
  staple = true
}
