//
//  NewAccountOrLogin.swift
//  soccer game
//
//  Created by Mark Schwartz on 2/13/18.
//  Copyright © 2018 MG 6. All rights reserved.
//
import SpriteKit
//import Alamofire

class MainMenu: SKScene {
    
    var title : SKLabelNode?
    var practiceBG : SKNode?
    var joinBG : SKNode?
    var comm : SKNode?
    
    override func didMove(to view: SKView) {
        print("got to main menu")
        
        //get scene subnodes
        self.back = self.childNode(withName: "Back Button")
        
        
        
        fadeInLabel(label: self.title)
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        //if necesarry nodes and one touch exist
        if let t = touches.first ,let back = self.back{
            let point = t.location(in: self)
            
            //see if touch contains first
            if back.contains(point){
                moveToMainMenu()
            }
        }
        
    }
    
    func fadeInLabel(label : SKLabelNode?){
        if let nonOptLabel = label{
            nonOptLabel.alpha = 0.0
            nonOptLabel.run(SKAction.fadeIn(withDuration: 2.0))
        }
    }
    
    
}
