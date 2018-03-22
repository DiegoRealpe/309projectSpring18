//
//  NewAccountOrLogin.swift
//  soccer game
//
//  Created by Mark Schwartz on 2/13/18.
//  Copyright © 2018 MG 6. All rights reserved.
//
import SpriteKit
import UIKit
import FacebookCore
import FacebookLogin
import FBSDKLoginKit
//import Alamofire



class NewAccountOrLogin: SKScene{
    
    
    var back : SKNode?
    var logout: SKNode?
    var viewController: UIViewController?
    var logOutLabel: SKLabelNode?
    var isLoggedIn = AccessToken.current != nil
    
    
    override func didMove(to view: SKView) {
        print("got to accounts")
        
        //get scene subnodes
        self.back = self.childNode(withName: "Back Button")
        self.logout = self.childNode(withName: "Logout")
        self.logOutLabel = logout?.childNode(withName: "LogoutLabel") as? SKLabelNode
        
         isLoggedIn = AccessToken.current != nil
        
        if(isLoggedIn)
        {
            logOutLabel?.text = "Logout"
        }
        else
        {
            logOutLabel?.text = "Login"
        }
        
       print(isLoggedIn)
    
        
        
        
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        
        //if necesarry nodes and one touch exist
        if let t = touches.first ,let back = self.back, let logout = self.logout{
            let point = t.location(in: self)
            
            //see if touch contains first
            if (back.contains(point)){
                self.moveToScene(.mainMenu)
            }
            else if (logout.contains(point))
            {
                if(isLoggedIn)//if we're logged in and want to logout
                {
                    let loginManager = LoginManager()
                    loginManager.logOut()
                    print("Logged out")
                    
                    logOutLabel?.text = "Login"
                
                }
                else
                {
                 /*   let fbLoginManager : FBSDKLoginManager = FBSDKLoginManager()
                    fbLoginManager.logIn(withReadPermissions: ["email"], from: self.viewController) { (result, error) -> Void in
                        if (error == nil){
                            let fbloginresult : FBSDKLoginManagerLoginResult = result!
                            if(fbloginresult.grantedPermissions.contains("email"))
                            {
                                
                            }
                        }
                    }
                    */
                    
                }
                
                
                
                
                
            }
        }
        
    }
    
    func fadeInLabel(label : SKLabelNode?){
        if let nonOptLabel = label{
            nonOptLabel.alpha = 0.0
            nonOptLabel.run(SKAction.fadeIn(withDuration: 2.0))
        }
    }
    
    func goToLogin(){
        var vc: UIViewController = UIViewController()
        vc = self.view!.window!.rootViewController!
        vc.performSegue(withIdentifier: "", sender: vc)
    }

    
}

