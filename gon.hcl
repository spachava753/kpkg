source = [
  "./release/darwin/kpkg"
]

bundle_id = "net.thespblog.kpkg"

apple_id {
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "8195B209BB79A5FF33C728DC1BBA3B5F30D5869D"
}

dmg {
  output_path = "./release/darwin/kpkg.dmg",
  volume_name = "kpkg"
}

zip {
  output_path = "./release/darwin/kpkg.zip"
}
