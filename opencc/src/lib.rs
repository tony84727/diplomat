use crate::generated_binding::{
    opencc_close, opencc_convert_utf8, opencc_error, opencc_open,
    OPENCC_DEFAULT_CONFIG_SIMP_TO_TRAD, OPENCC_DEFAULT_CONFIG_TRAD_TO_SIMP,
};
use std::convert::TryInto;
use std::ffi::CStr;

#[allow(unused)]
#[allow(non_upper_case_globals)]
#[allow(non_camel_case_types)]
#[allow(non_snake_case)]
mod generated_binding;

pub struct Convertor {
    opencc: generated_binding::opencc_t,
}

pub enum Mode {
    S2T,
    T2S,
}

#[derive(Debug, PartialEq)]
pub enum Error {
    Unknown(String),
}

impl Convertor {
    pub fn new(mode: Mode) -> Self {
        unsafe {
            Self {
                opencc: opencc_open(
                    CStr::from_bytes_with_nul(match mode {
                        Mode::S2T => OPENCC_DEFAULT_CONFIG_SIMP_TO_TRAD,
                        Mode::T2S => OPENCC_DEFAULT_CONFIG_TRAD_TO_SIMP,
                    })
                    .unwrap()
                    .as_ptr(),
                ),
            }
        }
    }

    pub fn convert(self, input: &str) -> Result<String, Error> {
        use std::ffi::CString;
        let cstring = CString::new(input).unwrap();
        let size = cstring.as_bytes_with_nul().len();
        let out;
        unsafe {
            let ptr = opencc_convert_utf8(self.opencc, cstring.as_ptr(), size.try_into().unwrap());
            if ptr.is_null() {
                let error = CStr::from_ptr(opencc_error());
                let error = error.to_str().unwrap_or_default();
                return Err(Error::Unknown(String::from(error)));
            }
            out = CString::from_raw(ptr).into_string().unwrap();
        }
        Ok(out)
    }
}
impl Drop for Convertor {
    fn drop(&mut self) {
        unsafe {
            opencc_close(self.opencc);
        }
    }
}
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_s2t() {
        let convertor = Convertor::new(Mode::T2S);
        assert_eq!(Ok(String::from("测试")), convertor.convert("測試"));
    }

    #[test]
    fn test_t2s() {
        let convertor = Convertor::new(Mode::S2T);
        assert_eq!(Ok(String::from("測試")), convertor.convert("测试"));
    }
}
