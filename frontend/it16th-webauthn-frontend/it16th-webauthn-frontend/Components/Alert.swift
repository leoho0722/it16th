//
//  Alert.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/21.
//

import UIKit

@MainActor
class Alert {
    
    class func showWith(title: String,
                        message: String,
                        confirmTitle: String,
                        confirm: (() -> Void)? = nil,
                        vc: UIViewController) {
        let alert = UIAlertController(title: title, message: message, preferredStyle: .alert)
        let confirmAction = UIAlertAction(title: confirmTitle, style: .default) { _ in
            confirm?()
        }
        alert.addAction(confirmAction)
        vc.present(alert, animated: true)
    }
}
