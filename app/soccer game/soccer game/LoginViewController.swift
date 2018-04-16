//
//  LoginViewController.swift
//  soccer game
//
//  Created by rtoepfer on 4/14/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import UIKit
import FBSDKLoginKit
import FacebookCore
import FacebookLogin
import Alamofire

class LoginViewController: UIViewController, FBSDKLoginButtonDelegate{
    
    @IBOutlet weak var ChooseNicknameView: UIView!
    @IBOutlet weak var nicknameInput: UITextField!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        if FBSDKAccessToken.current() != nil {
            let userToken = AccessToken.current.unsafelyUnwrapped.authenticationToken
                
            sendCRUDServiceLoginRequest(FBToken: userToken)
        }
        else{
            let loginButton = FBSDKLoginButton()
            loginButton.center = view.center
            loginButton.delegate = self // Remember to set the delegate of the loginButton
            view.addSubview(loginButton)
        }

    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    

    func loginButton(_ loginButton: FBSDKLoginButton!, didCompleteWith result: FBSDKLoginManagerLoginResult!, error: Error!)
    {
        
        if error != nil{
            // Process error
        }
        else if result.isCancelled {
            // Handle cancellations
        }
        else {
            
            sendCRUDServiceLoginRequest(FBToken : AccessToken.current!.authenticationToken)
            loginButton.removeFromSuperview()
        }
    }
    
    func loginButtonDidLogOut(_ loginButton: FBSDKLoginButton!) {
        let loginButton = FBSDKLoginButton()
        loginButton.center = view.center
        loginButton.delegate = self // Remember to set the delegate of the loginButton
        view.addSubview(loginButton)
    }
    
    //returns application token
    func sendCRUDServiceLoginRequest(FBToken : String) {
        
        let requestURL = "http://\(CommunicationProperties.crudServiceHost):\(CommunicationProperties.crudServicePort)/player/login"
        
        let headers: HTTPHeaders = [
            "FacebookToken": FBToken,
            ]
        
        print("\n")
        print("\n")
        
        print("logging in with CRUD Service at URL",requestURL,"\nwith headers : \(headers)")
        
        Alamofire.request(requestURL , method : .get , headers : headers)
            .responseString(completionHandler: loginRequestResponse(_:))
    }
    
    //function to do things with response from server to a login request by client
    func loginRequestResponse(_ response : DataResponse<String>)
    {
        
        print("🍆🍆🍆status code",response.response!.statusCode) //todo handle 404 without it being a fatal error
        
        
        if(response.response!.statusCode == 202)
        {
            self.moveToGameViewController()
        }
        else if(response.response!.statusCode == 404)//if player not found --> create account
        {
            print("lets create a new user 🐶")
            showCreateAccountView()
        }
        else if(response.response!.statusCode == 400)
        {
            print("SERVER MACHINE 🅱️ROKE")
        }
        
    }
    
    func showCreateAccountView(){
        self.ChooseNicknameView.isHidden = false
    }
    
    @IBAction func nicknameButtonPressed(){
        print("nickname submit button pressed")
        createAccount(FBToken: AccessToken.current!.authenticationToken, nickname: self.nicknameInput.text!)
    }
    
    
    func createAccount(FBToken : String,nickname : String)
    {
        print("\n")
        print("\n")
        
        let requestURL = "http://\(CommunicationProperties.crudServiceHost):\(CommunicationProperties.crudServicePort)/player/register"
        
        let headers: HTTPHeaders = [
            "FacebookToken": FBToken,
            ]
        
        
        let parameter:Parameters = ["Nickname":nickname]
        
        print("creating account with CRUD Service at URL",requestURL,"\nwith headers : \(headers)")
        
        Alamofire.request(requestURL, method: .post, parameters: parameter, encoding: JSONEncoding.default, headers: headers).responseString(completionHandler: createAccountResponse(_:))
        
        LocalPlayerInfo.username = nickname
    }
    func createAccountResponse(_ response : DataResponse<String>)
    {
        print("💌💌💌status code",response.response!.statusCode)
        
        if(response.response!.statusCode == 201)
        {
            moveToGameViewController()
        }
        else
        {
            print("Account creation machine 🅱️ROKE")
        }
    }
    
    func moveToGameViewController(){
        let storyboard = UIStoryboard(name: "Main", bundle: nil)
        let vc = storyboard.instantiateViewController(withIdentifier: "Game View Controller") as! GameViewController
        self.present(vc, animated: true, completion: nil)
    }

}
