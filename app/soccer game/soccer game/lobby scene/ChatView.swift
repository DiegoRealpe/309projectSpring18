//
//  ChatView.swift
//  soccer game
//
//  Created by rtoepfer on 4/2/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import UIKit

class ChatView: UIView {

    @IBOutlet var contentView: UIView!
    
    override init(frame: CGRect) {
        super.init(frame: frame)
        initCommon()
    }
    
    required init?(coder aDecoder: NSCoder) {
        super.init(coder: aDecoder)
        initCommon()
    }
    
    
    func initCommon(){
        Bundle.main.loadNibNamed("ChatView", owner: self, options: nil)
        addSubview(contentView)
        
        contentView.frame = CGRect(x: 0, y: 0, width: self.frame.width, height: self.frame.height)

    }

}
