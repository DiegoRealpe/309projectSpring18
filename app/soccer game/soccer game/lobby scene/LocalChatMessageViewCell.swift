//
//  TableViewCell.swift
//  chattest
//
//  Created by rtoepfer on 4/1/18.
//  Copyright © 2018 rtoepfer. All rights reserved.
//

import UIKit

class LocalChatMessageViewCell: UITableViewCell {
    
    @IBOutlet weak var messageLabel: UILabel!
    
    override func awakeFromNib() {
        super.awakeFromNib()
        // Initialization code
    }

    override func setSelected(_ selected: Bool, animated: Bool) {
        super.setSelected(selected, animated: animated)

        // Configure the view for the selected state
    }

}
