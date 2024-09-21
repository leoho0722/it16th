//
//  String+Extensions.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/21.
//

import Foundation

extension String {
    
    /// 將 Base64 編碼字串轉換為 Base64URL 編碼字串
    func base64EncodedToBase64URLEncoded() -> String {
        var base64URLEncoded = self
            .replacingOccurrences(of: "+", with: "-")
            .replacingOccurrences(of: "/", with: "_")
        if base64URLEncoded.count % 4 != 0 {
            base64URLEncoded.append(String(repeating: "=", count: 4 - base64URLEncoded.count % 4))
        }
        return base64URLEncoded
    }

    /// 將 Base64 編碼字串轉換為 Raw Base64URL 編碼字串
    func base64EncodedToBase64RawURLEncoded() -> String {
        return self
            .replacingOccurrences(of: "+", with: "-")
            .replacingOccurrences(of: "/", with: "_")
            .replacingOccurrences(of: "=", with: "")
    }
}
