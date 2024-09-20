//
//  PasskeysViewController.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/19.
//

import UIKit

class PasskeysViewController: UIViewController {
    
    // MARK: - IBOutlet
    
    @IBOutlet weak var txfUsername: UITextField!
    @IBOutlet weak var btnRegistration: UIButton!
    @IBOutlet weak var btnAuthentication: UIButton!
    
    // MARK: - Properties
    
    
    
    // MARK: - LifeCycle
    
    override func viewDidLoad() {
        super.viewDidLoad()
        setupUI()
    }
    
    override func viewWillAppear(_ animated: Bool) {
        super.viewWillAppear(animated)
    }
    
    override func viewIsAppearing(_ animated: Bool) {
        super.viewIsAppearing(animated)
    }
    
    override func viewWillLayoutSubviews() {
        super.viewWillLayoutSubviews()
    }
    
    override func viewDidLayoutSubviews() {
        super.viewDidLayoutSubviews()
    }
    
    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)
    }
    
    override func viewWillDisappear(_ animated: Bool) {
        super.viewWillDisappear(animated)
    }
    
    override func viewDidDisappear(_ animated: Bool) {
        super.viewDidDisappear(animated)
    }
    
    // MARK: - UI Settings
    
    fileprivate func setupUI() {
        // Username TextField
        txfUsername.placeholder = "請輸入使用者名稱"
        
        // Registration Button
        setupButton(btn: btnRegistration, title: "Registration", color: .systemPink, alpha: 0.2)
        
        // Authentication Button
        setupButton(btn: btnAuthentication, title: "Authentication", color: .systemBlue, alpha: 0.2)
    }
    
    fileprivate func setupButton(btn: UIButton, title: String, color: UIColor, alpha: CGFloat) {
        btn.tintColor = color.withAlphaComponent(alpha)
        btn.configuration?.baseForegroundColor = color
        btn.setTitle(title, for: .normal)
    }
    
    // MARK: - IBAction
    
    @IBAction func btnRegistrationClicked(_ sender: UIButton) {
        
    }
    
    @IBAction func btnAuthenticationClicked(_ sender: UIButton) {
        
    }
    
    // MARK: - Functions
    
    
}

// MARK: - Extensions



// MARK: - Protocol



// MARK: - Previews

#Preview {
    PasskeysViewController()
}
