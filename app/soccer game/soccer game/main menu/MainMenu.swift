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
    
    override func didMove(to view: SKView) {
        print("got to main menu")
        
        self.title = self.childNode(withName: "TitleLabel") as? SKLabelNode
        
        fadeInLabel(label: self.title)
        
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        let wrappedPractice = self.childNode(withName: "Practice Background")
        let wrappedJoin = self.childNode(withName: "Join Background")
        let wrappedComm = self.childNode(withName: "Comm Label")
        
        if let t = touches.first ,let practice = wrappedPractice, let join = wrappedJoin, let comm = wrappedComm{
            let point = t.location(in: self)
            if practice.contains(point){
                moveToGameScene()
                blinkLabel(label: self.title)
            }else if join.contains(point){
                moveToMatchMakingScene()
            }
            else if comm.contains(point){
                moveToCommTestScreen()
            }
        }
        
    }
    
    
    
    func fadeInLabel(label : SKLabelNode?){
        if let nonOptLabel = label{
            nonOptLabel.alpha = 0.0
            nonOptLabel.run(SKAction.fadeIn(withDuration: 2.0))
        }
    }
    
    func blinkLabel(label : SKLabelNode?){
        if let nonOptLabel = label{
            nonOptLabel.run(SKAction.init(named: "Pulse")!, withKey: "fadeInOut")
        }
    }
    
}
