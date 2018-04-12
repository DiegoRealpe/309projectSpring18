//
//  EmojiTextField.swift
//  soccer game
//
//  Created by rtoepfer on 4/11/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import UIKit

class EmojiTextField : UITextField {
    
    override var textInputMode: UITextInputMode? {
        for mode in UITextInputMode.activeInputModes {
            if mode.primaryLanguage == "emoji" {
                return mode
            }
        }
        return nil
    }
    
}
