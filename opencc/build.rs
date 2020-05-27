use std::path::PathBuf;

fn main() {
    println!("cargo:rustc-link-lib=opencc");
    println!("cargo:rerun-if-changed=wrapper.h");
    let bindings = bindgen::Builder::default()
        .header("wrapper.h")
        .parse_callbacks(Box::new(bindgen::CargoCallbacks))
        .generate()
        .expect("Unable to generate bindings");

    let src_path = PathBuf::from("src");
    bindings
        .write_to_file(src_path.join("generated_binding.rs"))
        .expect("Couldn't write bindings!");
}
