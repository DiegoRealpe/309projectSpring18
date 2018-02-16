//
//  MainMenu.swift
//  soccer game
//
//  Created by rtoepfer on 1/28/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit
import Alamofire

class MainMenu: SKScene {
    
    var title : SKLabelNode?
    var practiceBG : SKNode?
    var joinBG : SKNode?
    var comm : SKNode?
    var account:SKNode?
    override func didMove(to view: SKView) {
        print("got to main menu")
        
        //get scene subnodes
        self.title = self.childNode(withName: "TitleLabel") as? SKLabelNode
        self.practiceBG = self.childNode(withName: "Practice Background")
        self.joinBG = self.childNode(withName: "Join Background")
        self.comm = self.childNode(withName: "Comm Label")
        self.account = self.childNode(withName: "Account")
        
        fadeInLabel(label: self.title)
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        //if necesarry nodes and one touch exist
        if let t = touches.first ,let practice = self.practiceBG, let join = self.joinBG, let comm = self.comm, let account = self.account{
            let point = t.location(in: self)
            
            //see if touch contains first
            if practice.contains(point){
                moveToGameScene()
            }else if join.contains(point){
                moveToMatchMakingScene()
            }
            else if comm.contains(point){
                moveToCommTestScreen()
            }
            else if account.contains(point)
            {
                moveToAccountScreen()
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
