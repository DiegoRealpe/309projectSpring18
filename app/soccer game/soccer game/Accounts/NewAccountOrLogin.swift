//
//  NewAccountOrLogin.swift
//  soccer game
//
//  Created by Mark Schwartz on 2/13/18.
//  Copyright © 2018 MG 6. All rights reserved.
//
import SpriteKit
import UIKit

import FBSDKLoginKit
//import Alamofire

class NewAccountOrLogin: SKScene{
    
    var dict : [String : AnyObject]!
    var back : SKNode?
    var login: SKNode?
    
    
    override func didMove(to view: SKView) {
        print("got to accounts")
        
        //get scene subnodes
        self.back = self.childNode(withName: "Back Button")
        self.login = self.childNode(withName: "Login")
        
       
        
        
        
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        
        //if necesarry nodes and one touch exist
        if let t = touches.first ,let back = self.back, let login = self.login{
            let point = t.location(in: self)
            
            //see if touch contains first
            if back.contains(point){
                self.moveToScene(.mainMenu)
            }
            else if login.contains(point)
            {
                loginButtonClicked()
            }
        }
        
    }
    
    func fadeInLabel(label : SKLabelNode?){
        if let nonOptLabel = label{
            nonOptLabel.alpha = 0.0
            nonOptLabel.run(SKAction.fadeIn(withDuration: 2.0))
        }
    }
    
    @objc func loginButtonClicked() {
        let loginManager = LoginManager()
        loginManager.logIn([ .publicProfile ], viewController: self) { loginResult in
            switch loginResult {
            case .failed(let error):
                print(error)
            case .cancelled:
                print("User cancelled login.")
            case .success(let grantedPermissions, let declinedPermissions, let accessToken):
                self.getFBUserData()
            }
        }
    }
    
    //function is fetching the user data
    func getFBUserData(){
        if((FBSDKAccessToken.current()) != nil){
            FBSDKGraphRequest(graphPath: "me", parameters: ["fields": "id, name, picture.type(large), email"]).start(completionHandler: { (connection, result, error) -> Void in
                if (error == nil){
                    self.dict = result as! [String : AnyObject]
                    print(result!)
                    print(self.dict)
                }
            })
        }
    }
    

    
}

