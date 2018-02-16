//
//  SKSceneTransitionExtension.swift
//  soccer game
//
//  Created by rtoepfer on 1/29/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit

extension SKScene{
    
    func moveToGameScene(){
        
        if let scene = GameScene(fileNamed: "GameScene"){
            // Set the scale mode to scale to fit the window
            scene.scaleMode = .aspectFill
            
            
            // Present the scene
            if let view = self.view{
                view.presentScene(scene)
            }
        }
        
    }
    
    func moveToMatchMakingScene(){
        
        if let scene = GameScene(fileNamed: "MatchMakingScene"){
            // Set the scale mode to scale to fit the window
            scene.scaleMode = .aspectFill
            
            
            // Present the scene
            if let view = self.view{
                view.presentScene(scene)
            }
        }
        
    }
    
    func moveToMainMenu(){
        
        if let scene = GameScene(fileNamed: "MainMenu"){
            // Set the scale mode to scale to fit the window
            scene.scaleMode = .aspectFill
            
            
            // Present the scene
            if let view = self.view{
                view.presentScene(scene)
            }
        }
        
    }
    
    func moveToCommTestScreen(){
        
        if let scene = GameScene(fileNamed: "CommTestScreen"){
            // Set the scale mode to scale to fit the window
            scene.scaleMode = .aspectFill
            
            
            // Present the scene
            if let view = self.view{
                view.presentScene(scene)
            }
        }
        
    }
    func moveToAccountScreen(){
        
        if let scene = GameScene(fileNamed: "NewAccountOrLogin"){
            // Set the scale mode to scale to fit the window
            scene.scaleMode = .aspectFill
            
            
            // Present the scene
            if let view = self.view{
                view.presentScene(scene)
            }
        }
        
    }
    
   
    
    
}
