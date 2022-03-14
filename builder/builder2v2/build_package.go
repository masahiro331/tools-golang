// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package builder2v2

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/utils"
)

// BuildPackageSection2_2 creates an SPDX Package (version 2.2), returning
// that package or error if any is encountered. Arguments:
//   - packageName: name of package / directory
//   - dirRoot: path to directory to be analyzed
//   - pathsIgnore: slice of strings for filepaths to ignore
func BuildPackageSection2_2(packageName string, dirRoot string, pathsIgnore []string) (*spdx.Package2_2, error) {
	// build the file section first, so we'll have it available
	// for calculating the package verification code
	filepaths, err := utils.GetAllFilePaths(dirRoot, pathsIgnore)
	if err != nil {
		return nil, err
	}

	files := map[spdx.ElementID]*spdx.File2_2{}
	fileNumber := 0
	for _, fp := range filepaths {
		newFilePatch := filepath.FromSlash("./" + fp)
		newFile, err := BuildFileSection2_2(strings.Replace(newFilePatch, string(filepath.Separator), "/", -1), dirRoot, fileNumber)
		if err != nil {
			return nil, err
		}
		files[newFile.FileSPDXIdentifier] = newFile
		fileNumber++
	}

	// get the verification code
	code, err := utils.GetVerificationCode2_2(files, "")
	if err != nil {
		return nil, err
	}

	// now build the package section
	pkg := &spdx.Package2_2{
		PackageName:                 packageName,
		PackageSPDXIdentifier:       spdx.ElementID(fmt.Sprintf("Package-%s", packageName)),
		PackageDownloadLocation:     "NOASSERTION",
		FilesAnalyzed:               true,
		IsFilesAnalyzedTagPresent:   true,
		PackageVerificationCode:     code,
		PackageLicenseConcluded:     "NOASSERTION",
		PackageLicenseInfoFromFiles: []string{},
		PackageLicenseDeclared:      "NOASSERTION",
		PackageCopyrightText:        "NOASSERTION",
		Files:                       files,
	}

	return pkg, nil
}
