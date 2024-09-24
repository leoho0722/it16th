//
//  PasskeysViewController.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/19.
//

import AuthenticationServices
import SwiftHelpers
import UIKit

class PasskeysViewController: UIViewController {
    
    // MARK: - IBOutlet
    
    @IBOutlet weak var txfUsername: UITextField!
    @IBOutlet weak var txfDisplayName: UITextField!
    @IBOutlet weak var btnRegistration: UIButton!
    @IBOutlet weak var btnAuthentication: UIButton!
    
    // MARK: - Properties
    
    private let passkeysManager = PasskeysManager(domain: RequestConfiguration.Host.rpServer.rawValue)
    
    // MARK: - LifeCycle
    
    override func viewDidLoad() {
        super.viewDidLoad()
        setupUI()
        passkeysManager.delegate = self
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
        
        // DisplayName TextField
        txfDisplayName.placeholder = "請輸入使用者顯示名稱"
        
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
        guard let username = txfUsername.text else {
            Alert.showWith(title: "警告",
                           message: "請輸入使用者名稱！",
                           confirmTitle: "確認",
                           vc: self)
            return
        }
        guard let displayName = txfDisplayName.text else {
            Alert.showWith(title: "警告",
                           message: "請輸入使用者顯示名稱！",
                           confirmTitle: "確認",
                           vc: self)
            return
        }
        passkeysBeginRegistration(username: username, displayName: displayName)
    }
    
    @IBAction func btnAuthenticationClicked(_ sender: UIButton) {
        guard let username = txfUsername.text else {
            Alert.showWith(title: "警告",
                           message: "請輸入使用者名稱！",
                           confirmTitle: "確認",
                           vc: self)
            return
        }
        passkeysBeginAuthentication(username: username)
    }
    
    // MARK: - Functions
    
    private func pushToHome() {
        DispatchQueue.main.async {
            let homeVC = HomeViewController()
            self.navigationController?.pushViewController(homeVC, animated: true)
        }
    }
    
    private func passkeysBeginRegistration(username: String, displayName: String) {
        Task {
            do {
                let authenticatorSelection = AuthenticatorSelectionCriteria(authenticatorAttachment: .platform,
                                                                            residentKey: .preferred)
                let request = AttestationOptionsRequest(username: username,
                                                        displayName: displayName,
                                                        authenticatorSelection: authenticatorSelection,
                                                        attestation: .direct)
                let requestConfiguration = RequestConfiguration(method: .post,
                                                                scheme: .https,
                                                                host: .rpServer,
                                                                endpoint: .beginRegistration,
                                                                body: request)
                let response: AttestationOptionsResponse = try await NetworkManager.shared.request(with: requestConfiguration)
                
                guard let window = self.view.window else {
                    fatalError("The view was not in the app's view hierarchy!")
                }
                passkeysManager.registration(username: username, challenge: response.challenge, anchor: window)
            } catch {
                var errorMessage: Any
                switch error {
                case let networkError as NetworkError:
                    switch networkError {
                    case .badRequest(let response), .internalServerError(let response):
                        let decoder = JSONDecoder()
                        let decodedResponse = try! decoder.decode(CommonResponse.self, from: response)
                        errorMessage = decodedResponse.errorMessage
                    default:
                        errorMessage = error
                    }
                default:
                    errorMessage = error
                }
                Alert.showWith(title: "錯誤",
                               message: "WebAuthn 產生註冊資訊失敗！\n錯誤訊息為：\(errorMessage)",
                               confirmTitle: "確認",
                               vc: self)
            }
        }
    }
    
    private func passkeysFinishRegistration(clientDataJSON: Data,
                                            attestationObject: Data?,
                                            credentialID: Data) {
        Task {
            do {
                let base64RawURLEncodedClientDataJSON = clientDataJSON.base64EncodedString().base64EncodedToBase64RawURLEncoded()
                let base64RawURLEncodedAttestationObject = attestationObject?.base64EncodedString().base64EncodedToBase64RawURLEncoded()
                let base64RawURLEncodedCredentialID = credentialID.base64EncodedString().base64EncodedToBase64RawURLEncoded()
                
                let authenticatorAttestationResponse = AttestationResultsRequest.AuthenticatorAttestationResponse(clientDataJSON: base64RawURLEncodedClientDataJSON,
                                                                                                                  attestationObject: base64RawURLEncodedAttestationObject)
                let request = AttestationResultsRequest(id: base64RawURLEncodedCredentialID,
                                                        response: authenticatorAttestationResponse,
                                                        getClientExtensionResults: .init(),
                                                        type: .publicKey)
                let requestConfiguration = RequestConfiguration(method: .post,
                                                                scheme: .https,
                                                                host: .rpServer,
                                                                endpoint: .finishRegistration,
                                                                body: request)
                let response: CommonResponse = try await NetworkManager.shared.request(with: requestConfiguration)
                
                if response.status == "ok" {
                    Alert.showWith(title: "成功",
                                   message: "WebAuthn Registration 已完成！",
                                   confirmTitle: "確認",
                                   vc: self)
                }
            } catch {
                var errorMessage: String
                switch error {
                case let networkError as NetworkError:
                    switch networkError {
                    case .badRequest(let response), .internalServerError(let response):
                        let decoder = JSONDecoder()
                        let decodedResponse = try! decoder.decode(CommonResponse.self, from: response)
                        errorMessage = decodedResponse.errorMessage
                    default:
                        errorMessage = error.localizedDescription
                    }
                default:
                    errorMessage = error.localizedDescription
                }
                Alert.showWith(title: "錯誤",
                               message: "WebAuthn Registration 驗證註冊資訊失敗！\n錯誤訊息為：\(errorMessage)",
                               confirmTitle: "確認",
                               vc: self)
            }
        }
    }
    
    private func passkeysBeginAuthentication(username: String) {
        Task {
            do {
                let request = AssertionOptionsRequest(username: username, userVerification: .preferred)
                let requestConfiguration = RequestConfiguration(method: .post,
                                                                scheme: .https,
                                                                host: .rpServer,
                                                                endpoint: .beginAuthentication,
                                                                body: request)
                let response: AssertionOptionsResponse = try await NetworkManager.shared.request(with: requestConfiguration)
                
                guard let window = self.view.window else {
                    fatalError("The view was not in the app's view hierarchy!")
                }
                
                passkeysManager.authentication(challenge: response.challenge,
                                               anchor: window,
                                               preferImmediatelyAvailableCredentials: true)
            } catch {
                var errorMessage: Any
                switch error {
                case let networkError as NetworkError:
                    switch networkError {
                    case .badRequest(let response), .internalServerError(let response):
                        let decoder = JSONDecoder()
                        let decodedResponse = try! decoder.decode(CommonResponse.self, from: response)
                        errorMessage = decodedResponse.errorMessage
                    default:
                        errorMessage = error
                    }
                default:
                    errorMessage = error
                }
                Alert.showWith(title: "錯誤",
                               message: "WebAuthn Authentication 產生登入資訊失敗！\n錯誤訊息為：\(errorMessage)",
                               confirmTitle: "確認",
                               vc: self)
            }
        }
    }
    
    private func passkeysFinishAuthentication(clientDataJSON: Data,
                                              authenticatorData: Data?,
                                              signature: Data?,
                                              credentialID: Data,
                                              userID: Data?) {
        Task {
            do {
                let base64URLEncodedClientDataJSON = clientDataJSON.base64EncodedString().base64EncodedToBase64RawURLEncoded()
                let base64URLEncodedAuthenticatorData = authenticatorData?.base64EncodedString().base64EncodedToBase64RawURLEncoded()
                let base64URLEncodedSignature = signature?.base64EncodedString().base64EncodedToBase64RawURLEncoded()
                let base64URLEncodedCredentialID = credentialID.base64EncodedString().base64EncodedToBase64RawURLEncoded()
                let base64URLEncodedUserID = userID?.base64EncodedString().base64EncodedToBase64RawURLEncoded()
                
                let authenticatorAssertionResponse = AssertionResultsRequest.AuthenticatorAssertionResponse(authenticatorData: base64URLEncodedAuthenticatorData,
                                                                                                            signature: base64URLEncodedSignature,
                                                                                                            userHandle: base64URLEncodedUserID,
                                                                                                            clientDataJSON: base64URLEncodedClientDataJSON)
                
                let request = AssertionResultsRequest(id: base64URLEncodedCredentialID,
                                                      response: authenticatorAssertionResponse,
                                                      getClientExtensionResults: .init(),
                                                      type: .publicKey)
                let requestConfiguration = RequestConfiguration(method: .post,
                                                                scheme: .https,
                                                                host: .rpServer,
                                                                endpoint: .finishAuthentication,
                                                                body: request)
                let response: CommonResponse = try await NetworkManager.shared.request(with: requestConfiguration)
                
                if response.status == "ok" {
                    Alert.showWith(title: "成功",
                                   message: "WebAuthn Authentication 已完成！",
                                   confirmTitle: "確認",
                                   confirm: pushToHome,
                                   vc: self)
                }
            } catch {
                var errorMessage: Any
                switch error {
                case let networkError as NetworkError:
                    switch networkError {
                    case .badRequest(let response), .internalServerError(let response):
                        let decoder = JSONDecoder()
                        let decodedResponse = try! decoder.decode(CommonResponse.self, from: response)
                        errorMessage = decodedResponse.errorMessage
                    default:
                        errorMessage = error
                    }
                default:
                    errorMessage = error
                }
                Alert.showWith(title: "錯誤",
                               message: "WebAuthn Authentication 驗證登入資訊失敗！\n錯誤訊息為：\(errorMessage)",
                               confirmTitle: "確認",
                               vc: self)
            }
        }
    }
}

// MARK: - Extensions

// MARK: PasskeysManagerDelegate

extension PasskeysViewController: @preconcurrency PasskeysManagerDelegate {
    
    func passkeysManager(controller: ASAuthorizationController, didCompleteWithError error: any Error) {
        Alert.showWith(title: "錯誤",
                       message: error.localizedDescription,
                       confirmTitle: "確認",
                       vc: self)
    }
}

// MARK: PasskeysRegistration

extension PasskeysViewController: @preconcurrency PasskeysRegistration {
    
    func passkeysManager(with credentialRegistration: ASAuthorizationPlatformPublicKeyCredentialRegistration) {
        let clientDataJSON = credentialRegistration.rawClientDataJSON
        let attestationObject = credentialRegistration.rawAttestationObject
        let credentialID = credentialRegistration.credentialID
        let attachment = credentialRegistration.attachment
        
        passkeysFinishRegistration(clientDataJSON: clientDataJSON,
                                   attestationObject: attestationObject,
                                   credentialID: credentialID)
    }
}

// MARK: PasskeysAuthentication

extension PasskeysViewController: @preconcurrency PasskeysAuthentication {
    
    func passkeysManager(with credentialAssertion: ASAuthorizationPlatformPublicKeyCredentialAssertion) {
        let clientDataJSON = credentialAssertion.rawClientDataJSON
        let authenticatorData = credentialAssertion.rawAuthenticatorData
        let signature = credentialAssertion.signature
        let userID = credentialAssertion.userID
        let credentialID = credentialAssertion.credentialID
        let attachment = credentialAssertion.attachment
        
        passkeysFinishAuthentication(clientDataJSON: clientDataJSON,
                                     authenticatorData: authenticatorData,
                                     signature: signature,
                                     credentialID: credentialID,
                                     userID: userID)
    }
}

// MARK: - Protocol



// MARK: - Previews

//#Preview {
//    PasskeysViewController()
//}
